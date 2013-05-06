package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"
)

type ActivityIssueComment struct {
	IssueComments map[int]IssueCommentPayload
}

type Comment struct {
	Body string
	User GithubUser
}

func (c Comment) ShortBody() string {
	msg := strings.Split(c.Body, "\n")[0]
	if len(msg) > 80 {
		return msg[0:77] + "..."
	}
	return msg
}

type IssueCommentPayload struct {
	Payload struct {
		Action  string
		Issue   Issue // in issue.go
		Comment Comment
	}
}

const long_issue_comment_template = `{{ range $num, $payload := .IssueComments }}    #{{$num}}: {{$payload.Payload.Issue.Title}}
        {{ $payload.Payload.Comment.User.Login }} commented on issue {{$num}}: {{ $payload.Payload.Comment.ShortBody }}
{{end}}`

const short_issue_comment_template = `Comments on {{ range $num, $payload := .IssueComments }}#{{$num}} {{end}}
`

func issueCommentRender(activities []Activity, long_template bool) string {
	var metas = make(map[int]IssueCommentPayload, len(activities))
	for _, activity := range activities {
		var payload IssueCommentPayload
		err := json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil {
			log.Print("Error decoding issue comment meta for pk:%s: %s", activity.Id, err)
		}

		if payload.Payload.Issue.Number < 1 {
			log.Print("Malformed Issue Comment Payload: %s", activity.Meta)
			continue
		}

		metas[payload.Payload.Issue.Number] = payload
	}

	template_input := ActivityIssueComment{metas}
	tmpl := template.New("IssueCommentFragment")

	if long_template {
		_, err := tmpl.Parse(long_issue_comment_template)
		if err != nil {
			fmt.Println("Error with activity fragment parsing. ", err)
		}
	} else {
		_, err := tmpl.Parse(short_issue_comment_template)
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
