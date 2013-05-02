package main

import (
	"fmt"
	"encoding/json"
	"bytes"
	"text/template"
)

type CreatePayload struct {
	Payload CreateMeta
}

type CreateMeta struct {
	Ref_type string
	Ref string
	Master_branch string
	Description string
}

type ActivityCreate struct {
	Created []CreateMeta
}

const long_create_template = `
{{range .Created}}    Created {{.Ref_type}} {{.Ref}}
{{end}}
`

const short_create_template = `
    Created {{len .Created}} branches/refs.
`

func create_render(activities []Activity, long_template bool) string {
	var metas = make([]CreateMeta, len(activities))
	for i, activity := range activities {
		var payload CreatePayload
		err :=json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil { fmt.Println("Error decoding meta: ", err) }

		metas[i] = payload.Payload
	}
	
	template_input := ActivityCreate{metas}
	tmpl := template.New("CreateFragment")

	if long_template {
		_, err := tmpl.Parse(long_create_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	} else {
		_, err := tmpl.Parse(short_create_template)
		if err != nil { fmt.Println("Error with activity fragment parsing. ", err) }
	}
	
	var b bytes.Buffer
	err := tmpl.Execute(&b, template_input)
	if err != nil { fmt.Println("Error with activity rendering. ", err) }
	return b.String()
}
