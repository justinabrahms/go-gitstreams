package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func initTestDb(t *testing.T) *sql.DB {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_TEST_DB")
	if len(db_name) == 0 {
		t.Fatal("Must specify a db so we don't accidentally overwrite anything.")
	}
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, db_name))
	if err != nil {
		t.Fatal("Couldn't connect to DB.", err)
	}
	return db
}

func make_user(id int, username, email string, db *sql.DB, t *testing.T) {
	_, err := db.Exec(`INSERT INTO auth_user VALUES (
                   $1, $2,
                   'first', 'last', 
                   $3, 'pass', 
                   TRUE, TRUE, TRUE, now(), now()
                 )`, id, username, email)
	if err != nil {
		t.Fatal("Unable to create user.", err)
	}
}

func make_repo(id int, user, project string, db *sql.DB, t *testing.T) {
	_, err := db.Exec(`INSERT INTO streamer_repo VALUES (
                   $1, $2, $3, NULL, NULL
                 )`, id, user, project)
	if err != nil {
		t.Fatal("Unable to create repo.", err)
	}
}

func make_userprofile(id, uid int, time_interval string, db *sql.DB, t *testing.T) {
	_, err := db.Exec(`INSERT INTO streamer_userprofile
                 (id, user_id, max_time_interval_between_emails) VALUES (
                   $1, $2, $3
                 )`, id, uid, time_interval)
	if err != nil {
		t.Fatal("Unable to create userprofile.", err)
	}
}

func make_userprofile_repo(id, userprofile_id, repo_id int, db *sql.DB, t *testing.T) {
	_, err := db.Exec(`INSERT INTO streamer_userprofile_repos
                 (id, userprofile_id, repo_id) VALUES (
                   $1, $2, $3
                 )`, id, userprofile_id, repo_id)
	if err != nil {
		t.Fatal("Unable to create userprofile_repo.", err)
	}
}

func TestGetUser_Exists(t *testing.T) {
	db := initTestDb(t)
	dbc := DbController{db}
	defer dbc.Close()

	defer db.Exec("DELETE from auth_user")
	make_user(1, "username", "em@a.il", db, t)

	u, err := dbc.GetUser(1)
	if err != nil {
		t.Fatal(err)
	}

	if u.Id != 1 {
		t.Fatal("Incorrect user id.")
		return
	}
	if u.username != "username" {
		t.Fatal("Incorrect username.")
		return
	}
	if u.Email != "em@a.il" {
		t.Fatal("Incorrect email.")
		return
	}
}

func TestGetUser_NoExists(t *testing.T) {
	db := initTestDb(t)
	dbc := DbController{db}
	defer dbc.Close()

	_, err := dbc.GetUser(1)
	if err == nil {
		t.Fatal("No error when accessing invalid user id.")
		return
	}
}

func TestGetUserRepos_None(t *testing.T) {
	db := initTestDb(t)
	dbc := DbController{db}
	defer dbc.Close()
	uid := 1

	defer db.Exec("DELETE FROM auth_user;")
	make_user(uid, "username", "em@a.il", db, t)
	defer db.Exec("DELETE FROM streamer_userprofile;")
	make_userprofile(2, uid, "D", db, t) // Daily, should really be enum or something.
	defer db.Exec("DELETE FROM streamer_repo;")
	make_repo(3, "user", "repo", db, t)
	defer db.Exec("DELETE FROM streamer_userprofile_repos;")
	make_userprofile_repo(4, 2, 3, db, t)

	ur, err := dbc.GetUserRepos(uid)
	if err != nil {
		t.Fatal("Can't get user's repos.", err)
		return
	}
	if len(ur) != 1 {
		fmt.Println(ur)
		t.Fatal("Got user's repos, but expected len 1, but got ", len(ur))
		return
	}
}

func testGetUserRepos_One(t *testing.T) {
}

func testGetUserRepos_Many(t *testing.T) {
}

func testGetRepoActivity_OtherRepo(t *testing.T) {

}

func testGetRepoActivity_OtherUser(t *testing.T) {

}

func testGetRepoActivity_Old(t *testing.T) {

}

func testGetRepoActivity_AlreadySent(t *testing.T) {

}

func testGetRepoActivity_NeverSent(t *testing.T) {

}

func testMarkUserRepoSent(t *testing.T) {

}
