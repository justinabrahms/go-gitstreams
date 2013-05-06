package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

type ActivityPush struct {
	Pushes       []PushMeta
	TotalCommits int
}

type PushPayload struct {
	Payload PushMeta
}

type PushMeta struct {
	Commits []Commit
	Ref     string // eg refs/heads/master
	Head    string // head SHA
}

const long_push_template = `{{ range .Pushes }}{{range .Commits}}    {{.ShortSha}} {{ .Author.Name }} -- {{ .ShortCommit }}
{{ end }}
{{ end }}`

// {{$commit.ShortSha}}
const short_push_template = `{{ .TotalCommits }} commits. {{range .Pushes}}{{range .Commits}} {{.ShortSha}} {{end}}
{{end}}
`

func pushRender(activities []Activity, long_template bool) string {
	var metas = make([]PushMeta, len(activities))
	var total_commits = 0
	for i, activity := range activities {
		var payload PushPayload
		err := json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil {
			fmt.Println("Error decoding meta: ", err)
		}

		metas[i] = payload.Payload
		total_commits += len(payload.Payload.Commits)
	}

	template_input := ActivityPush{metas, total_commits}
	tmpl := template.New("ActivityFragment")

	if long_template {
		_, err := tmpl.Parse(long_push_template)
		if err != nil {
			fmt.Println("Error with activity fragment parsing. ", err)
		}
	} else {
		_, err := tmpl.Parse(short_push_template)
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
