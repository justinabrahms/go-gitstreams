package main

import (
	"fmt"
	"os"
	"bytes"
	"time"
	"encoding/json"
	"text/template"
	"database/sql"
	"strings"
	_ "github.com/go-sql-driver/mysql"
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

type Activity struct {
	id int
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

type ActivityPush struct {
	Pushes []PushMeta
	TotalCommits int
}

type PushPayload struct {
	 Payload PushMeta
}

type PushMeta struct {
	Commits []Commit
	Ref string // eg refs/heads/master
	Head string // head SHA
}

const long_push_template = `{{ range .Pushes }}{{range .Commits}}    {{.ShortSha}} {{ .Author.Name }} -- {{ .ShortCommit }}
{{ end }}
{{ end }}`

const short_push_template = `{{ .TotalCommits }} commits. 
{{range $index, $p := .Pushes}}{{range $p.Commits}}{{if $index}}, {{end}}{{.ShortSha}}{{end}}{{end}}
`

type ActivityDelete struct {
	Deleted []DeleteMeta
}

type DeletePayload struct {
	Payload DeleteMeta
}

type DeleteMeta struct {
	Ref_type string
	Ref string
}

const long_delete_template = `{{range .Deleted}}    Deleted {{.Ref_type}} {{.Ref}}{{end}}
`
const short_delete_template = `{{len .Deleted}} deleted branches/refs.`


type ActivityPullRequest struct {
	PullRequests map[int]PullRequest // number -> pull request
}

type PullRequestPayload struct {
	Payload PullRequestMeta
}

type PullRequestMeta struct {
	Number int
	Action string
	Pull_request PullRequest
}

type PullRequest struct {
	Number int
	State string // enum?
	Title string
	Body string
	Head Treeish
	Base Treeish

	// These are in the PR, but aren't any reasons to capture it.
	
	// Merged_by GithubUser

	// Created_at time.Time
	// Updated_at time.Time
	// Closed_at time.Time
	// Merged_at time.Time

	// Comments int
	// Commits int
	// Additions int
	// Deletions int
	// Changed_files int
}

const long_pr_template = `{{range $num, $pr := .PullRequests}}
    PR:{{.Number}} {{.Head.User.Login}} -- {{.Title}}
{{end}}
`

const short_pr_template = `{{len .PullRequests}} pull requests.
{{range $num, $pr := .PullRequests}}{{$num}}{{end}}
`

type CreatePayload struct {
	Payload CreateMeta
}

type CreateMeta struct {
	Ref_type string
	Ref string
	Master_branch string
	Description string
}

type ActivityCreate struct {
	Created []CreateMeta
}

const long_create_template = `
{{range .Created}}    Created {{.Ref_type}} {{.Ref}}
{{end}}
`

const short_create_template = `
    Created {{len .Created}} branches/refs.
`

type WatchPayload struct {
	Actor GithubUser
	Repo GithubApiRepo
	Payload WatchMeta
}

type WatchMeta struct {
	Action string
}

type ActivityWatch struct {
	Watched []WatchPayload
}

const long_watch_template = `
{{range .Watched}}    {{.Actor.Login}} {{.Payload.Action}} watching {{.Repo.Name}}
{{end}}
`

const short_watch_template = `
    {{len .Watched}} watch events.
`


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

func pull_request_render(activities []Activity, long_template bool) string {
	var metas = make(map[int]PullRequest, len(activities))
	for _, activity := range activities {
		var payload PullRequestPayload
		err :=json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil { fmt.Println("Error decoding meta: ", err) }
		metas[payload.Payload.Number] = payload.Payload.Pull_request
	}

	template_input := ActivityPullRequest{metas}
	tmpl := template.New("PullRequestFragment")

	if long_template {
		_, err := tmpl.Parse(long_pr_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	} else {
		_, err := tmpl.Parse(short_pr_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	}
	

	var b bytes.Buffer
	err := tmpl.Execute(&b, template_input)
	if err != nil { fmt.Println("Error with activity rendering. ", err) }
	
	return b.String()
}

func delete_render(activities []Activity, long_template bool) string { 
	var metas = make([]DeleteMeta, len(activities))
	for i, activity := range activities {
		var payload DeletePayload
		err :=json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil { fmt.Println("Error decoding meta: ", err) }

		metas[i] = payload.Payload
	}
	
	template_input := ActivityDelete{metas}
	tmpl := template.New("DeleteFragment")

	if long_template {
		_, err := tmpl.Parse(long_delete_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	} else {
		_, err := tmpl.Parse(short_delete_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	}
	
	var b bytes.Buffer
	err := tmpl.Execute(&b, template_input)
	if err != nil { fmt.Println("Error with activity rendering. ", err) }
	return b.String()
}

func create_render(activities []Activity, long_template bool) string {
	var metas = make([]CreateMeta, len(activities))
	for i, activity := range activities {
		var payload CreatePayload
		err :=json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil { fmt.Println("Error decoding meta: ", err) }

		metas[i] = payload.Payload
	}
	
	template_input := ActivityCreate{metas}
	tmpl := template.New("CreateFragment")

	if long_template {
		_, err := tmpl.Parse(long_create_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	} else {
		_, err := tmpl.Parse(short_create_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	}
	
	var b bytes.Buffer
	err := tmpl.Execute(&b, template_input)
	if err != nil { fmt.Println("Error with activity rendering. ", err) }
	return b.String()
}

func watch_render(activities []Activity, long_template bool) string {
	var metas = make([]WatchPayload, len(activities))
	for i, activity := range activities {
		var payload WatchPayload
		err :=json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil { fmt.Println("Error decoding meta: ", err) }
		metas[i] = payload
	}
	
	template_input := ActivityWatch{metas}
	tmpl := template.New("WatchFragment")

	if long_template {
		_, err := tmpl.Parse(long_watch_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	} else {
		_, err := tmpl.Parse(short_watch_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	}
	
	var b bytes.Buffer
	err := tmpl.Execute(&b, template_input)
	if err != nil { fmt.Println("Error with activity rendering. ", err) }
	return b.String()
}

func issue_comment_render(activities []Activity, long_template bool) string { return "" }
func fork_render(activities []Activity, long_template bool) string { return "" }
func issue_render(activities []Activity, long_template bool) string { return "" }
func gist_render(activities []Activity, long_template bool) string { return "" }
func follow_render(activities []Activity, long_template bool) string { return "" }
func commit_comment_render(activities []Activity, long_template bool) string { return "" }
func pull_request_comment_render(activities []Activity, long_template bool) string { return "" }
func wiki_render(activities []Activity, long_template bool) string { return "" }
func member_render(activities []Activity, long_template bool) string { return "" }
func public_render(activities []Activity, long_template bool) string { return "" }
func download_render(activities []Activity, long_template bool) string { return "" }
func fork_apply_render(activities []Activity, long_template bool) string { return "" }
func team_add_render(activities []Activity, long_template bool) string { return "" }

// Handles taking a list of Push-type activities and returning a formatted template
func push_render(activities []Activity, long_template bool) string {
	var metas = make([]PushMeta, len(activities))
	var total_commits = 0
	for i, activity := range activities {
		var payload PushPayload
		err :=json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil { fmt.Println("Error decoding meta: ", err) }

		metas[i] = payload.Payload
		total_commits += len(payload.Payload.Commits)
	}
	
	template_input := ActivityPush{metas, total_commits}
	tmpl := template.New("ActivityFragment")

	if long_template {
		_, err := tmpl.Parse(long_push_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	} else {
		_, err := tmpl.Parse(short_push_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	}
	
	var b bytes.Buffer
	err := tmpl.Execute(&b, template_input)
	if err != nil { fmt.Println("Error with activity rendering. ", err) }
	return b.String()
}

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
			response += fn(activities, true)
		} else {
			fmt.Println("Not sure how to render activites of type ", activity_type)
		}
	}
	return response
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_DB")))
	if (err != nil) {
		fmt.Println("Unable to connect to mysql. ", err)
		os.Exit(1)
	}
	defer db.Close()

	// var user_id = 1;
	// repos, err := get_users_repos(db, user_id)
	repo := GithubRepo{43567, "", ""}
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
	}

	// build map of repo -> activities
	// go routine w/ channel of activities, 
	//   throw down activities on per-repo basis
	//   collect them when done and collate.
	response := repo_to_template(repo, activities, activity_type_to_renderer)
	fmt.Println(response)

	return
}
