package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
)

var token = flag.String("oauth_token", "", "OAuth token of user")

func getEmailsForToken(token string) (emails []github.UserEmail, err error) {
	t := &oauth.Transport{
		Token: &oauth.Token{
			AccessToken: token,
		},
	}
	client := github.NewClient(t.Client())
	emails, _, err = client.Users.ListEmails()
	return
}

func main() {
	flag.Parse()

	if *token == "" {
		flag.Usage()
		return
	}

	emails, err := getEmailsForToken(*token)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Printf("Emails: %v\n", emails)
}
