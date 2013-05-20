package main

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"log"
	"strconv"
	"strings"
	"time"
)

func NewDbController(user, pass, database string) (c DbController, err error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=localhost sslmode=disable", user, pass, database))
	if err == nil {
		c = DbController{db}
	}
	return
}

type DbController struct {
	db *sql.DB
}

func (d *DbController) Close() {
	d.db.Close()
}

func (d *DbController) GetUser(uid int) (u User, err error) {
	row := d.db.QueryRow(
		`SELECT id, username, email
		   FROM auth_user
		   WHERE id = $1;`, uid)
	err = row.Scan(&u.Id, &u.username, &u.Email)
	return
}

func (d *DbController) GetUsers() (users []User, err error) {
	rows, err := d.db.Query(
		`SELECT id, username, email
		   FROM auth_user;`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(&u.Id, &u.username, &u.Email)
		if err != nil {
			return
		}
		users = append(users, u)
	}
	return
}

func (d *DbController) GetUserRepos(user_id int) ([]GithubRepo, error) {
	// guess as to initial size. Likely very few users following < 10 repos.
	var repos = make([]GithubRepo, 10)

	// TODO(justinabrahms): Should really alter the schema such
	// that the join from user <-> repo doesn't go through
	// userprofiles.
	rows, err := d.db.Query(
		`SELECT r.id, username, project_name
	           FROM streamer_repo r
		   JOIN streamer_userprofile_repos upr ON upr.repo_id = r.id 
		   JOIN streamer_userprofile up ON up.id = upr.userprofile_id
                   WHERE up.user_id = $1;`, user_id)
	// what's a good way to handle this not working?
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var username, project_name string
	var pk int
	for rows.Next() {
		err = rows.Scan(&pk, &username, &project_name)
		if err != nil {
			log.Print("Error fetching user's repos: ", user_id, err)
		}
		var repo = GithubRepo{pk, username, project_name}
		repos = append(repos, repo)
	}

	return repos, err
}

func (d *DbController) GetRepoActivity(repo *GithubRepo, userId int) (activity_list []Activity, err error) {
	activity_list = make([]Activity, 0)
	rows, err := d.db.Query(
		`SELECT a.id, a.event_id, a.type, a.created_at, ghu.name, r.username, r.project_name, meta
		  FROM streamer_activity a
		  JOIN streamer_repo r on r.id=a.repo_id
		  JOIN streamer_githubuser ghu on ghu.id=a.user_id
		  JOIN streamer_userprofile_repos upr on r.id=upr.repo_id
                  JOIN streamer_userprofile up on up.id=upr.userprofile_id
		  WHERE r.id = $1
                  AND up.user_id = $2
		  AND a.created_at > (NOW() - INTERVAL '5 days') -- don't send things more than a few days old. Think, new users who subscribe to rails/rails
		  AND (upr.last_sent is null -- hasn't been sent at all
		    OR a.created_at > upr.last_sent); --  or hasn't been sent since we've gotten new stuff`,
		repo.Id, userId)
	if err != nil {
		return
	}
	defer rows.Close()

	var (
		pk, github_id, repo_id                                 int
		activity_type, username, repo_project, repo_user, meta string
		created_at                                             time.Time
	)

	for rows.Next() {
		err = rows.Scan(&pk, &github_id, &activity_type, &created_at, &username, &repo_user, &repo_project, &meta)
		if err != nil {
			log.Print("Error getting activity for repo: ", repo, err)
		}

		activity_list = append(activity_list,
			Activity{pk, github_id, activity_type, created_at, username,
				GithubRepo{repo_id, repo_user, repo_project}, meta})
	}

	return
}

func (d *DbController) MarkUserRepoSent(user User, repos []GithubRepo) (err error) {
	// finds the streamer_userprofile_repo row for the repo / user
	// combo, mark its last_sent as now
	ids := make([]string, 0)
	for _, repo := range repos {
		// TODO: why are they 0?
		if repo.Id != 0 {
			ids = append(ids, strconv.FormatInt(int64(repo.Id), 10))
		}
	}

	// There is likely a better way to get parameterization, but
	// it wasn't working for me with ?'s.
	_, err = d.db.Exec(
		`UPDATE streamer_userprofile_repos
		   SET last_sent=NOW()
                   WHERE userprofile_id = (
                     SELECT id FROM streamer_userprofile where id=$1)
		   AND repo_id IN ($2)`, user.Id, strings.Join(ids, ","))
	return
}
