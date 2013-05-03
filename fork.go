package main

import (
	"log"
	"fmt"
	"encoding/json"
	"bytes"
	"text/template"
)


type ActivityFork struct {
	Forks []ForkPayload
}

type ForkPayload struct {
	Payload struct {
		Forkee GithubApiRepo
	}
	Repo GithubApiRepo
	Actor GithubUser
}

const long_fork_template = `{{range .Forks}}    {{.Actor.Login}} forked to {{.Payload.Forkee.FullName}}
{{end}}`

// TODO: maybe add the fork destinations?
const short_fork_template = `{{len .Forks}} new forks
`

func fork_render(activities []Activity, long_template bool) string {
	var metas = make([]ForkPayload, len(activities))
	for i, activity := range activities {
		var payload ForkPayload
		err :=json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil { log.Print("Error decoding Fork meta for pk:%s: ", activity.Id, err) }

		metas[i] = payload
	}

	
	template_input := ActivityFork{metas}
	tmpl := template.New("ForkFragment")

	if long_template {
		_, err := tmpl.Parse(long_fork_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	} else {
		_, err := tmpl.Parse(short_fork_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	}
	
	var b bytes.Buffer
	err := tmpl.Execute(&b, template_input)
	if err != nil { fmt.Println("Error with activity rendering. ", err) }
	return b.String()
}
