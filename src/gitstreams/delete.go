package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"
)

type ActivityDelete struct {
	Deleted []DeleteMeta
}

type DeletePayload struct {
	Payload DeleteMeta
}

type DeleteMeta struct {
	Ref_type string
	Ref      string
}

const long_delete_template = `{{range .Deleted}}    Deleted {{.Ref_type}} {{.Ref}}
{{end}}`
const short_delete_template = `{{len .Deleted}} deleted branches/refs.`

func delete_render(activities []Activity, long_template bool) string {
	var metas = make([]DeleteMeta, len(activities))
	for i, activity := range activities {
		var payload DeletePayload
		err := json.Unmarshal([]byte(activity.Meta), &payload)
		if err != nil {
			log.Print("Error decoding delete meta for pk:%s: %s", activity.Id, err)
		}

		metas[i] = payload.Payload
	}

	template_input := ActivityDelete{metas}
	tmpl := template.New("DeleteFragment")

	if long_template {
		_, err := tmpl.Parse(long_delete_template)
		if err != nil {
			fmt.Println("Error with activity fragment parsing. ", err)
		}
	} else {
		_, err := tmpl.Parse(short_delete_template)
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
