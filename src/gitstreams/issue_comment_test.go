package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func build_icp(body string) (icp IssueCommentPayload, err error) {
	json_string := []byte(fmt.Sprintf(`{"payload": {"action": "test", "issue": {}, "comment":{"body": "%s"}}}`, body))
	err = json.Unmarshal(json_string, &icp)
	return
}

func assertEqual(t *testing.T, actual, expected string) {
	if expected != actual {
		t.Errorf("Expected body to be `%s`, not `%s`", expected, actual)
	}
}

func TestBodyShort(t *testing.T) {
	body := "short one"
	icp, err := build_icp(body)
	if err != nil {
		t.Errorf("%s", err)
	}
	assertEqual(t, icp.Payload.Comment.ShortBody(), body)
}

func TestBodyLong(t *testing.T) {
	body := "longer than 80 characater string which should be split before the complete end of the string else the test will fail. Really."
	icp, err := build_icp(body)
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(icp.Payload.Comment.ShortBody()) > 80 {
		t.Errorf("Didn't truncate it to 80 characters")
	}
	if !strings.HasSuffix(icp.Payload.Comment.ShortBody(), "...") {
		t.Errorf("Didn't add elipsis at the end.")
	}
}
