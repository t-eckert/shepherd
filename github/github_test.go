package github

import (
	"fmt"
	"testing"
)

func TestExecIssueList(t *testing.T) {
	repo := "t-eckert/shepherd"

	response, err := execIssueList(repo)

	if err != nil {
		fmt.Println(response)
		fmt.Println(err)
		t.Fail()
	}
}

func TestParseIssuesToList(t *testing.T) {
	issues := []byte("[{\"number\":2,\"title\":\"Issue 2\"},{\"number\":1,\"title\":\"Issue 1\"}]")

	expected := []Issue{
		{Number: 2, Title: "Issue 2"},
		{Number: 1, Title: "Issue 1"},
	}

	actual, _ := parseIssuesToList(issues)

	if len(expected) != len(actual) {
		fmt.Printf("given: %v\nexpected: %v\nactual: %v", issues, expected, actual)
		t.Fail()
	}
	for i, expectedIssue := range expected {
		if expectedIssue != actual[i] {
			fmt.Printf("given: %v\nexpected: %v\nactual: %v", issues, expectedIssue, actual[i])
			t.Fail()
		}
	}
}

func TestEditPrepend(t *testing.T) {
	params := []struct {
		givenIssueTitle string
		givenPrepend    string
		expected        string
	}{
		{"one:Issue title", "two", "two:Issue title"},
		{"Issue title", "two", "two:Issue title"},
	}

	for _, param := range params {
		actual := editPrepend(param.givenIssueTitle, param.givenPrepend)

		if param.expected != actual {
			fmt.Printf("given: %v, %v\nexpected: %v\nactual: %v\n", param.givenIssueTitle, param.givenPrepend, param.expected, actual)
			t.Fail()
		}
	}
}
