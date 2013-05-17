package main

import (
	"fmt"
	"log"
)

func repoToTemplate(repo GithubRepo, activities []Activity, render_map map[string]func([]Activity, bool) string) (response string) {
	if len(activities) == 0 {
		return ""
	}
	var activity_map = make(map[string][]Activity)
	for _, activity := range activities {
		// This seems like a lot of juggling. Is there a better way?
		arr := activity_map[activity.activity_type]
		if arr == nil {
			arr = make([]Activity, 0)
			activity_map[activity.activity_type] = arr
		}
		arr = append(arr, activity)
		activity_map[activity.activity_type] = arr
	}

	// activity_map: activity_type => []activity
	for activity_type, activities := range activity_map {
		fn, ok := render_map[activity_type]
		if ok {
			response += fn(activities, true)
		} else {
			log.Print("Not sure how to render activites of type ", activity_type)
		}
	}
	if len(response) > 0 {
		response = fmt.Sprintf("\n\n%s:\n%s", repo.FullName(), response)
	}
	return

}

func repoToString(db *DbController, repo GithubRepo, userId int, response chan string) {
	activities, err := db.GetRepoActivity(&repo, userId)
	if err != nil {
		log.Print("ERR: ", err)
	}

	activity_type_to_renderer := map[string]func([]Activity, bool) string{
		"P":  pushRender,
		"PR": pullRequestRender,
		"D":  deleteRender,
		"C":  createRender,
		"IC": issueCommentRender,
		"Gl": wikiRender, // Gl is for Gollum, Github's wiki thing.
		"I":  issueRender,
		"Pb": publicRender,

		// Watches and forks for repos are WAYY too spammy.
		// "W":  watchRender,
		// "F":  forkRender,
	}

	response <- repoToTemplate(repo, activities, activity_type_to_renderer)
}
