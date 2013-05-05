package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	_ "github.com/go-sql-driver/mysql"
	mailgun "github.com/riobard/go-mailgun"
)

type User struct {
	Id int
	username string
	Email string
}

type GithubUser struct {
	Id int
	Login string
	Url string
}

// db thing.
type GithubRepo struct {
	Id int
	User string // should probably be Login
	RepoName string // should probably be Name
}

func (g *GithubRepo) FullName() string {
	return fmt.Sprintf("%s/%s", g.User, g.RepoName)
}

type GithubApiRepo struct {
	Id int
	Owner GithubUser
	Name NString
	Description NString
}

func (g *GithubApiRepo) FullName() string {
	return fmt.Sprintf("%s/%s", g.Owner.Login, g.Name)
}

type Activity struct {
	Id int
	github_id int
	activity_type string // should be enum
	created_at time.Time
	Username string
	repo GithubRepo // this isn't going to be memory efficient
	Meta string // full payload of json object
}

type Treeish struct {
	Label string
	Sha string
	Repo GithubApiRepo
	Html_url string
	User GithubUser
}

type Commit struct {
	Sha string
	Message string
	Author struct {
		Name string
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

// Nullable String
type NString string

func (n *NString) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(n))
}


func get_users_repos(db *sql.DB, user_id int) ([]GithubRepo, error) {
	var repo_count int

	// TODO(justinabrahms): Should really alter the schema such
	// that the join from user <-> repo doesn't go through
	// userprofiles.
	row := db.QueryRow("SELECT count(*) FROM streamer_repo r" +
	" JOIN streamer_userprofile_repos upr ON upr.repo_id = r.id " +
	" JOIN streamer_userprofile up ON up.id = upr.userprofile_id" +
	" WHERE up.user_id = ?;", user_id)

	err := row.Scan(&repo_count)
	if (err != nil) {
		return nil, err
	}
	var repos = make([]GithubRepo, repo_count)

	rows, err := db.Query("SELECT r.id, username, project_name FROM streamer_repo r" +
	" JOIN streamer_userprofile_repos upr ON upr.repo_id = r.id " +
	" JOIN streamer_userprofile up ON up.id = upr.userprofile_id" +
	" WHERE up.user_id = ?;", user_id)
	// what's a good way to handle this not working?
	if (err != nil) {
		return nil, err
	}
	defer rows.Close()

	var username, project_name string
	var pk int
	for rows.Next() {
		err = rows.Scan(&pk, &username, &project_name)
		if (err != nil) {
			fmt.Println("Error!: ", err)
		}
		var repo = GithubRepo{pk, username, project_name}
		repos = append(repos, repo)
	}

	return repos, err
}



func get_user(db *sql.DB, uid int) (u User, err error) {
	row := db.QueryRow(
		"SELECT id, username, email" +
		" FROM auth_user" +
		" WHERE id = ?", uid)
	err = row.Scan(&u.Id, &u.username, &u.Email)
	return
}

func get_repo_activity(db *sql.DB, repo *GithubRepo) (activity_list []Activity, err error){
	activity_list = make([]Activity, 0)
	rows, err := db.Query(
		"SELECT a.id, a.event_id, a.type, a.created_at, ghu.name, r.username, r.project_name, meta" +
		" FROM streamer_activity a" +
		" JOIN streamer_repo r on r.id=a.repo_id" +
		" JOIN streamer_githubuser ghu on ghu.id=a.user_id" +
		" JOIN streamer_userprofile_repos upr on r.id=upr.repo_id" +
		" WHERE r.id = ?" +
		" AND a.created_at > DATE_SUB(NOW(), INTERVAL 5 day)" + // don't send things more than a few days old. Think, new users who subscribe to rails/rails
		" AND (upr.last_sent is null" + // hasn't been sent at all
		"   OR a.created_at > upr.last_sent)",  // or hasn't been sent since we've gotten new stuff
		repo.Id)
	if (err != nil) {
		return
	}
	defer rows.Close()

	
	var (
		pk, github_id, repo_id int
		activity_type, username, repo_project, repo_user, meta, created_str string
		created_at time.Time
	)
	
	for rows.Next() {
		err = rows.Scan(&pk, &github_id, &activity_type, &created_str, &username, &repo_user, &repo_project, &meta)
		if (err != nil) {
			fmt.Println("ERR: ", err)
		}
		
		created_at, err = time.Parse("2006-01-02 15:04:05", created_str)
		if (err != nil) {
			fmt.Println("Can't parse the format of ", created_str)
			return nil, err
		}

		activity_list = append(activity_list, 
			Activity{pk, github_id, activity_type, created_at, username,
			GithubRepo{repo_id, repo_user, repo_project}, meta})
	}
	
	return
}

func commit_comment_render(activities []Activity, long_template bool) string { return "" }
func pull_request_comment_render(activities []Activity, long_template bool) string { return "" }
func member_render(activities []Activity, long_template bool) string { return "" }
func download_render(activities []Activity, long_template bool) string { return "" }
func fork_apply_render(activities []Activity, long_template bool) string { return "" }
func team_add_render(activities []Activity, long_template bool) string { return "" }
func gist_render(activities []Activity, long_template bool) string { return "" }
func follow_render(activities []Activity, long_template bool) string { return "" }

func repo_to_template(repo GithubRepo, activities []Activity, render_map map[string]func([]Activity, bool)string) string {
	if len(activities) == 0 {
		return ""
	}
	var activity_map = make(map[string][]Activity)
	for _, activity := range activities {
		// This seems like a lot of juggling. Is there a better way?
		arr := activity_map[activity.activity_type]
		if (arr == nil) {
			arr = make([]Activity, 0)
			activity_map[activity.activity_type] = arr
		}
		arr = append(arr, activity)
		activity_map[activity.activity_type] = arr
	}

	// activity_map: activity_type => []activity
	var response = fmt.Sprintf("%s:\n", repo.FullName())
	for activity_type, activities := range activity_map {
		fn, ok := render_map[activity_type]
		if ok {
			response += fn(activities, true)
		} else {
			log.Print("Not sure how to render activites of type ", activity_type)
		}
	}
	return response
}

func repo_to_string(db *sql.DB, repo GithubRepo, response chan string) {
	activities, err := get_repo_activity(db, &repo)
	if (err != nil) {
		fmt.Println("ERR: ", err)
		os.Exit(1)
	}

	activity_type_to_renderer := map[string]func([]Activity, bool)string {
		"P": push_render,
		"PR": pull_request_render,
		"D": delete_render,
		"C": create_render,
		"W": watch_render,
		"F": fork_render,
		"IC": issue_comment_render,
		"Gl": wiki_render, // Gl is for Gollum, Github's wiki thing.
		"I": issue_render,
		"Pb": public_render,
		// TODO: TA (team add), FA (fork apply),
		// TODO: DO (Download), RC (Pull Request Review Comment),
		// TODO: Fl (Follow), G (Gist)
	}

	response <- repo_to_template(repo, activities, activity_type_to_renderer)
}

func mark_user_repo_sent(db *sql.DB, user User, repos []GithubRepo) (err error) {
	// should find the streamer_userprofile_repo row for the repo
	// / user combo, mark its last_sent as now
	ids := make([]string, 0)
	for _, repo := range repos {
		// why are they 0?
		if repo.Id != 0 {
			ids = append(ids, strconv.FormatInt(int64(repo.Id), 10))
		}
	}

	// There is likely a better way to get parameterization, but
	// it wasn't working for me with ?'s.
	str := fmt.Sprintf(
		"UPDATE streamer_userprofile_repos " +
		"SET last_sent=NOW() " +
		"WHERE repo_id IN (%s)", strings.Join(ids, ","))
	_, err = db.Exec(str)
	return
}


var user_id = flag.Int("user_id", 1, "ID of user to output.")
var email_to = flag.String("email_to", "", "Email address of who should get the report.")
var mark_read = flag.Bool("mark_read", true, "Whether to mark activity sent as read. False means subsequent calls will send the same info.")

// TODO: Github Users

func main() {
	flag.Parse()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_DB")))
	if (err != nil) {
		log.Fatalf("Unable to connect to mysql. ", err)
	}
	defer db.Close()

	repos, err := get_users_repos(db, *user_id)
	if err != nil { 
		log.Fatalf("Error fetching user's repos. %s", err)
	}

	response_chan := make(chan string, len(repos))
	for _, repo := range repos {
		repo_to_string(db, repo, response_chan)
	}

	var response string
	for num_responses := 0; num_responses < len(repos); {
		select {
		case r := <-response_chan:
			num_responses++
			response += r
		}
	}

	// should fetch this user based on uid.
	user, err := get_user(db, *user_id)
	if err != nil { log.Fatalf("Couldn't return user %d. %s", *user_id, err) }
	if *mark_read {
		fmt.Println("Marking read.")
		err = mark_user_repo_sent(db, user, repos)
		if err != nil {
			log.Fatalf("Error updating repositories as sent.")
		}
	}

	if len(response) == 0 {
		log.Fatal("No content to email.")
	}

	if len(*email_to) > 0 {
		mg := mailgun.Open(os.Getenv("MAILGUN_API_KEY"))
		e := &Email{
			from: "justin@gitstreams.mailgun.org",
			to: []string{user.Email},
			subject: "Gitstreams Digest Email",
			text: response,
		}

		id, err := mg.Send(e)

		if err != nil {
			log.Fatal("Unable to send email. ", err)
		}
		log.Printf("MessageId = %s for uid:%s", id, user.Id)
	} else {
		fmt.Println("Would have sent the following email.")
		fmt.Println(response)
	}
	return
}
