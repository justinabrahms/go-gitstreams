package main

import (
	"encoding/json"
	"flag"
	"fmt"
	mailgun "github.com/riobard/go-mailgun"
	"log"
	"os"
	"strings"
	"time"
)

var user_id = flag.Int("user_id", 0, "ID of user to output.")
var all_users = flag.Bool("all_users", false, "Whether to email all users.")
var send_email = flag.Bool("send_email", false, "Whether to send the email to the user.")
var mark_read = flag.Bool("mark_read", true, "Whether to mark activity sent as read. False means subsequent calls will send the same info.")

type User struct {
	Id       int
	username string
	Email    string
}

type GithubUser struct {
	Id    int
	Login string
	Url   string
}

// db thing.
type GithubRepo struct {
	Id       int
	User     string // should probably be Login
	RepoName string // should probably be Name
}

func (g *GithubRepo) FullName() string {
	return fmt.Sprintf("%s/%s", g.User, g.RepoName)
}

type GithubApiRepo struct {
	Id          int
	Owner       GithubUser
	Name        NString
	Description NString
}

func (g *GithubApiRepo) FullName() string {
	return fmt.Sprintf("%s/%s", g.Owner.Login, g.Name)
}

type Activity struct {
	Id            int
	github_id     int
	activity_type string // should be enum
	created_at    time.Time
	Username      string
	repo          GithubRepo // this isn't going to be memory efficient
	Meta          string     // full payload of json object
}

type Treeish struct {
	Label    string
	Sha      string
	Repo     GithubApiRepo
	Html_url string
	User     GithubUser
}

type Commit struct {
	Sha     string
	Message string
	Author  struct {
		Name  string
		Email string
	}
}

func (c *Commit) ShortSha() string {
	return c.Sha[0:6]
}

func (c *Commit) ShortCommit() string {
	msg := strings.Split(c.Message, "\n")[0]
	if len(msg) > 80 {
		return msg[0:77] + "..."
	}
	return msg
}

// Nullable String, per https://groups.google.com/d/msg/golang-nuts/JOFWAqrTbUs/KZIGEPclSpwJ
type NString string

func (n *NString) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(n))
}

// TODO: Github Users

// TODO: Need to finish the following activity types: 
// - TA (team add)
// - FA (fork apply) 
// - DO (Download) 
// - RC (Pull Request Review Comment)
// - Fl (Follow) 
// - G (Gist)
// - M (Member)
// - CC (Commit Comment)

func main() {
	// Expects the following environment variables:
	// DB_USER, DB_PASS, DB_DB  (database user, password and database)
	// MAILGUN_API_KEY  (api key from mailgun.org)
	flag.Parse()

	c, err := NewDbController(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_DB"))
	if err != nil {
		log.Fatalf("Unable to connect to mysql. ", err)
	}
	defer c.Close()

	var users []User
	if *user_id != 0 {
		users = make([]User, 1)
		user, err := c.GetUser(*user_id)
		if err != nil {
			log.Fatalf("Couldn't return user %d. %s", *user_id, err)
		}
		users[0] = user
	} else if *all_users {
		users, err = c.GetUsers()
		if err != nil {
			log.Fatal("Couldn't fetch all users. %s", err)
		}
	} else {
		log.Fatal("You must specify either to mail all users or a specific user id.")
	}

	for _, user := range users {
		repos, err := c.GetUserRepos(user.Id)
		if err != nil {
			log.Print("Error fetching user's repos. %s", err)
			continue
		}

		response_chan := make(chan string, len(repos))
		for _, repo := range repos {
			repoToString(&c, repo, user.Id, response_chan)
		}

		var response string
		for num_responses := 0; num_responses < len(repos); {
			select {
			case r := <-response_chan:
				num_responses++
				response += r
			}
		}

		if *mark_read {
			fmt.Println("Marking read.")
			err = c.MarkUserRepoSent(user, repos)
			if err != nil {
				log.Print("Error updating repositories as sent.")
				continue
			}
		}

		if len(response) == 0 {
			log.Print("No content to email.")
			continue
		}

		if *send_email {
			mg := mailgun.Open(os.Getenv("MAILGUN_API_KEY"))
			e := &Email{
				from:    "justin@gitstreams.mailgun.org",
				to:      []string{user.Email},
				subject: "Gitstreams Digest Email",
				text:    response,
			}

			id, err := mg.Send(e)

			if err != nil {
				log.Print("Unable to send email. ", err)
			}
			log.Printf("MessageId = %s for uid:%s", id, user.Id)
		} else {
			fmt.Println("Would have sent the following email.")
			fmt.Println(response)
		}
	}

	return
}
