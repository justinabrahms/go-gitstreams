package main

import (
	"fmt"
	"os"
	"time"
	"encoding/json"
	"text/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	username string
}

type GithubRepo struct {
	pk int
	User string
	RepoName string
}

// JSON thing from GitHub
type Payload interface {}


type Activity struct {
	id int
	github_id int
	activity_type string // should be enum
	created_at time.Time
	Username string
	repo GithubRepo // this isn't going to be memory efficient
	Meta Payload // full payload of json object
}

type PushPayload struct {
	 Payload PushMeta
}

type PushMeta struct {
	Commits []Commit
	Ref string // eg refs/heads/master
	Head string // head SHA
}

type Commit struct {
	Sha string
	Message string
	Author string `json:"author.name"` // name
}

type ActivityTemplateInput struct {
	Repo GithubRepo
	Input map[string][]Activity
}

const activity_template = `
Repo: {{.Repo.User}}/{{.Repo.RepoName}}
{{ range .Input.P }}
  Push by {{ .Username }} -- {{range .Meta.Commits}}{{.Sha}}{{ end }}
{{ end }}`

const activity_entry = `activity entry
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

		// deserialize
		var payload Payload
		if (activity_type == "P") {
				var intermediate PushPayload
			err = json.Unmarshal([]byte(meta), &intermediate)
			// fmt.Println("Payload: ", payload)
			// var pushMeta PushMeta
			// err := json.Unmarshal([]byte(meta), &pushMeta)
			// if (err != nil) {
			// 	fmt.Println("Unable to decode push event: ", err)
			// }
			fmt.Println("Got a: ", payload, " from a: ", meta)
				payload = intermediate
			// payload = pushMeta
		}
		
		activity_list = append(activity_list, 
			Activity{pk, github_id, activity_type, created_at, username,
			GithubRepo{repo_id, repo_user, repo_project}, payload})
	}
	
	return
}

func repo_to_template(repo GithubRepo, activities []Activity) string {
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

	template_input := ActivityTemplateInput{repo, activity_map}

	tmpl, err := template.New("ActivityFragment").Parse(activity_template)
	if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	err = tmpl.Execute(os.Stdout, template_input)
	if err != nil { fmt.Println("Error with activity rendering. ", err) }
	return ""
}

func main() {
	fmt.Println("DB User: ", os.Getenv("DB_USER"))
	fmt.Println("DB Database: ", os.Getenv("DB_DB"))
	
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_DB")))
	if (err != nil) {
		fmt.Println("Unable to connect to mysql. ", err)
		os.Exit(1)
	}
	defer db.Close()

	// var user_id = 1;
	// repos, err := get_users_repos(db, user_id)
	repo := GithubRepo{43563, "exitio", "schematics"}
	activities, err := get_repo_activity(db, &repo)
	if (err != nil) {
		fmt.Println("ERR: ", err)
		os.Exit(1)
	}

	// build map of repo -> activities
	// go routine w/ channel of activities, 
	//   throw down activities on per-repo basis
	//   collect them when done and collate.
	repo_to_template(repo, activities)

	return
}
