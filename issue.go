package main

import (
	"log"
	"fmt"
	"encoding/json"
	"bytes"
	"text/template"
)


type ActivityIssue struct {
	Issues []IssuePayload
}

type Issue struct {
	State string // enum?
	Title string
	Body NString
	Number int
}

type IssueMeta struct {
	Action string
	Issue Issue
}

type IssuePayload struct {
	Payload IssueMeta
	Actor GithubUser
}


const long_issue_template = `{{range .Issues}}{{.Actor.Login}} {{.Payload.Action}} #{{.Payload.Issue.Number}}: {{.Payload.Issue.Title}}
{{end}}
`

const short_issue_template = `{{len .Issues}} new issues: {{ range $i, $issue := .Issues }}{{if $i}}, {{end}}#{{$issue.Payload.Issue.Number}}{{end}}
`

func issue_render(activities []Activity, long_template bool) string {
	var metas = make([]IssuePayload, len(activities))
	for i, activity := range activities {
		var payload IssuePayload
		err :=json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil { log.Print("Error decoding Issue meta for pk:%d: ", activity.Id, err) }
		metas[i] = payload
	}

	template_input := ActivityIssue{metas}
	tmpl := template.New("IssueFragment")

	if long_template {
		_, err := tmpl.Parse(long_issue_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	} else {
		_, err := tmpl.Parse(short_issue_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	}
	
	var b bytes.Buffer
	err := tmpl.Execute(&b, template_input)
	if err != nil { fmt.Println("Error with activity rendering. ", err) }
	return b.String()
}
