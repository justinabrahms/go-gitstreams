package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

type ActivityWiki struct {
	Wikis []WikiPayload
}

type WikiPayload struct {
	Payload WikiMeta
	Actor   GithubUser
}

type Page struct {
	Page_name string
	Title     string
	Action    string
	Sha       string
	Html_url  string
}

type WikiMeta struct {
	Pages []Page
}

const long_wiki_template = `{{ range .Wikis }}    {{ .Actor.Login }} {{ range $i, $page := .Payload.Pages}}{{if $i}}, {{end}}{{ $page.Action }} {{$page.Page_name}}{{end}}
{{end}}
`

const short_wiki_template = `{{ range .Wikis }}{{with .Payload}}    Updated pages: {{range $index, $page := .Pages}}{{if $index}}, {{end}}{{ $page.Title }}{{end}}{{end}}{{end}}
`

func wikiRender(activities []Activity, long_template bool) string {
	var metas = make([]WikiPayload, len(activities))
	for i, activity := range activities {
		var payload WikiPayload
		err := json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil {
			fmt.Println("Error decoding meta: ", err)
		}

		metas[i] = payload
	}

	template_input := ActivityWiki{metas}
	tmpl := template.New("WikiFragment")

	if long_template {
		_, err := tmpl.Parse(long_wiki_template)
		if err != nil {
			fmt.Println("Error with activity fragment parsing. ", err)
		}
	} else {
		_, err := tmpl.Parse(short_wiki_template)
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
