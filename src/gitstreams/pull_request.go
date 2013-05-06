package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

type ActivityPullRequest struct {
	PullRequests map[int]PullRequest // number -> pull request
}

type PullRequestPayload struct {
	Payload PullRequestMeta
}

type PullRequestMeta struct {
	Number       int
	Action       string
	Pull_request PullRequest
}

type PullRequest struct {
	Number int
	State  string // enum?
	Title  string
	Body   string
	Head   Treeish
	Base   Treeish

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
{{range $num, $pr := .PullRequests}}{{$num}} {{end}}
`

func pullRequestRender(activities []Activity, long_template bool) string {
	var metas = make(map[int]PullRequest, len(activities))
	for _, activity := range activities {
		var payload PullRequestPayload
		err := json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil {
			fmt.Println("Error decoding meta: ", err)
		}
		metas[payload.Payload.Number] = payload.Payload.Pull_request
	}

	template_input := ActivityPullRequest{metas}
	tmpl := template.New("PullRequestFragment")

	if long_template {
		_, err := tmpl.Parse(long_pr_template)
		if err != nil {
			fmt.Println("Error with activity fragment parsing. ", err)
		}
	} else {
		_, err := tmpl.Parse(short_pr_template)
		if err != nil {
			fmt.Println("Error with activity fragment parsing. ", err)
		}
	}

	var b bytes.Buffer
	err := tmpl.Execute(&b, template_input)
	if err != nil {
		fmt.Println("Error with activity rendering. ", err)
	}

	return b.String()
}
