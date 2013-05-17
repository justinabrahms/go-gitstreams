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

func TestGetUser_Exists(t *testing.T) {
	db := initTestDb(t)
	dbc := DbController{db}
	defer dbc.Close()

	_, err := db.Exec(`INSERT INTO auth_user VALUES (
                   1, 'username', 
                   'first', 'last', 
                   'em@a.il', 'pass', 
                   TRUE, TRUE, TRUE, now(), now()
                 )`)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Exec("DELETE from auth_user")

	u, err := dbc.GetUser(1)
	if err != nil {
		t.Fatal(err)
	}

	if u.Id != 1 {
		t.Fatal("Incorrect user id.")
	}
	if u.username != "username" {
		t.Fatal("Incorrect username.")
	}
	if u.Email != "em@a.il" {
		t.Fatal("Incorrect email.")
	}
}

func testGetUser_NoExists(t *testing.T) {

}

func testGetUserRepos_None(t *testing.T) {

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
