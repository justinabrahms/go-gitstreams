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


func repoToTemplate(repo GithubRepo, activities []Activity, render_map map[string]func([]Activity, bool) string) (response string) {
	if len(activities) == 0 {
		return ""
	}
	var activity_map = make(map[string][]Activity)
	for _, activity := range activities {
		// This seems like a lot of juggling. Is there a better way?
		arr := activity_map[activity.activity_type]
		if arr == nil {
			arr = make([]Activity, 0)
			activity_map[activity.activity_type] = arr
		}
		arr = append(arr, activity)
		activity_map[activity.activity_type] = arr
	}

	// activity_map: activity_type => []activity
	for activity_type, activities := range activity_map {
		fn, ok := render_map[activity_type]
		if ok {
			response += fn(activities, true)
		} else {
			log.Print("Not sure how to render activites of type ", activity_type)
		}
	}
	if len(response) > 0 {
		response = fmt.Sprintf("\n\n%s:\n%s", repo.FullName(), response)
	}
	return
	
}

func repoToString(db *sql.DB, repo GithubRepo, response chan string) {
	activities, err := getRepoActivity(db, &repo)
	if err != nil {
		fmt.Println("ERR: ", err)
		os.Exit(1)
	}

	activity_type_to_renderer := map[string]func([]Activity, bool) string{
		"P":  pushRender,
		"PR": pullRequestRender,
		"D":  deleteRender,
		"C":  createRender,
		"IC": issueCommentRender,
		"Gl": wikiRender, // Gl is for Gollum, Github's wiki thing.
		"I":  issueRender,
		"Pb": publicRender,

		// Watches and forks for repos are WAYY too spammy.
		// "W":  watchRender,
		// "F":  forkRender,
	}

	response <- repoToTemplate(repo, activities, activity_type_to_renderer)
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
		user, err := c.getUser(*user_id)
		if err != nil {
			log.Fatalf("Couldn't return user %d. %s", *user_id, err)
		}
		users[0] = user
	} else if *all_users {
		users, err = c.getUsers()
		if err != nil {
			log.Fatal("Couldn't fetch all users. %s", err)
		}
	} else {
		log.Fatal("You must specify either to mail all users or a specific user id.")
	}

	for _, user := range users {
		repos, err := c.getUserRepos(user.Id)
		if err != nil {
			log.Print("Error fetching user's repos. %s", err)
			continue
		}

		response_chan := make(chan string, len(repos))
		for _, repo := range repos {
			repoToString(db, repo, response_chan)
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
			err = c.markUserRepoSent(user, repos)
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
