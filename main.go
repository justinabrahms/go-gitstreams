package main

import (
	"fmt"
	"os"
	"time"
	"database/sql"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"flag"
)

type User struct {
	username string
}

type GithubUser struct {
	Id int
	Login string
	Url string
}

// db thing.
type GithubRepo struct {
	pk int
	User string // should probably be Login
	RepoName string // should probably be Name
}

type GithubApiRepo struct {
	Id int
	Owner GithubUser
	Name string
	Description string
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


// get user's followed users
// get user's followed repos
// for each ^, pull out activity (goroutine)
//   return rendered template bits
// join template bits


// Does not take over ownership of db.
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


func get_repo_activity(db *sql.DB, repo *GithubRepo) (activity_list []Activity, err error){
	rows, err := db.Query(
		"SELECT a.id, a.event_id, a.type, a.created_at, ghu.name, r.username, r.project_name, meta FROM streamer_activity a" +
		" JOIN streamer_repo r on r.id=a.repo_id" +
		" JOIN streamer_githubuser ghu on ghu.id=a.user_id" +
		" WHERE repo_id = ?", repo.pk)
	if (err != nil) {
		return nil, err
	}
	defer rows.Close()

	activity_list = make([]Activity, 0)
	
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

func gist_render(activities []Activity, long_template bool) string { return "" }
func follow_render(activities []Activity, long_template bool) string { return "" }
func commit_comment_render(activities []Activity, long_template bool) string { return "" }
func pull_request_comment_render(activities []Activity, long_template bool) string { return "" }
func member_render(activities []Activity, long_template bool) string { return "" }
func public_render(activities []Activity, long_template bool) string { return "" }
func download_render(activities []Activity, long_template bool) string { return "" }
func fork_apply_render(activities []Activity, long_template bool) string { return "" }
func team_add_render(activities []Activity, long_template bool) string { return "" }

func repo_to_template(repo GithubRepo, activities []Activity, render_map map[string]func([]Activity, bool)string) string {
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
	var response = ""
	for activity_type, activities := range activity_map {
		fn, ok := render_map[activity_type]
		if ok {
			response += fn(activities, false)
		} else {
			fmt.Println("Not sure how to render activites of type ", activity_type)
		}
	}
	return response
}

var repo_id = flag.Int("repo_id", 43567, "ID of github repo to dump.")

func main() {
	flag.Parse()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_DB")))
	if (err != nil) {
		fmt.Println("Unable to connect to mysql. ", err)
		os.Exit(1)
	}
	defer db.Close()

	// var user_id = 1;
	// repos, err := get_users_repos(db, user_id)
	repo := GithubRepo{*repo_id, "", ""}
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
	}

	// build map of repo -> activities
	// go routine w/ channel of activities, 
	//   throw down activities on per-repo basis
	//   collect them when done and collate.
	response := repo_to_template(repo, activities, activity_type_to_renderer)
	fmt.Println(response)

	return
}
