package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

type WatchPayload struct {
	Actor   GithubUser
	Repo    GithubApiRepo
	Payload WatchMeta
}

type WatchMeta struct {
	Action string
}

type ActivityWatch struct {
	Watched []WatchPayload
}

const long_watch_template = `
{{range .Watched}}    {{.Actor.Login}} {{.Payload.Action}} watching {{.Repo.Name}}
{{end}}
`

const short_watch_template = `
    {{len .Watched}} watch events.
`

func watchRender(activities []Activity, long_template bool) string {
	var metas = make([]WatchPayload, len(activities))
	for i, activity := range activities {
		var payload WatchPayload
		err := json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil {
			fmt.Println("Error decoding meta: ", err)
		}
		metas[i] = payload
	}

	template_input := ActivityWatch{metas}
	tmpl := template.New("WatchFragment")

	if long_template {
		_, err := tmpl.Parse(long_watch_template)
		if err != nil {
			fmt.Println("Error with activity fragment parsing. ", err)
		}
	} else {
		_, err := tmpl.Parse(short_watch_template)
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
