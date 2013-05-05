package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

type PublicTemplateInput struct {
	Payloads []PublicPayload
}

type PublicPayload struct {
	Actor GithubUser
	Repo  struct {
		Name string
	}
}

const long_public_template = `{{range .Payloads}}{{.Repo.Name}} was open sourced!
{{end}}
`

const short_public_template = long_public_template

func public_render(activities []Activity, long_template bool) string {
	var metas = make([]PublicPayload, len(activities))
	for i, activity := range activities {
		var payload PublicPayload
		err := json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil {
			fmt.Println("Error decoding meta: ", err)
		}

		metas[i] = payload
	}

	template_input := PublicTemplateInput{metas}
	tmpl := template.New("Public Fragment")

	if long_template {
		_, err := tmpl.Parse(long_public_template)
		if err != nil {
			fmt.Println("Error with activity fragment parsing. ", err)
		}
	} else {
		_, err := tmpl.Parse(short_public_template)
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
