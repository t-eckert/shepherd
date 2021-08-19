package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Issue struct {
	Number int
	Title  string
}

// FetchIssues uses the GitHub command line application to fetch all open issues from the given repository
func FetchIssues(repo string) ([]Issue, error) {
	output, err := execIssueList(repo)
	if err != nil {
		return nil, fmt.Errorf("could not fetch issues for %s: %v", repo, err)
	}

	issueList, err := parseIssuesToList(output)
	if err != nil {
		return nil, fmt.Errorf("could not parse issues for %s: %v", repo, err)
	}

	return issueList, nil
}

// ModifyPrepend iterates over the given issues and changes their title's to
// replace or add the prepended string
func ModifyPrepend(issues []Issue, repo, prepend string) error {
	for _, issue := range issues {
		newTitle := editPrepend(issue.Title, prepend)
		if err := execEditIssueTitle(repo, issue.Number, newTitle); err != nil {
			return fmt.Errorf("could not rename issue %d to %s", issue.Number, newTitle)
		}
	}

	return nil
}

// execIssueList uses the GitHub CLI to return a list of issues
func execIssueList(repo string) ([]byte, error) {
	output, err := exec.Command("gh", "issue", "list", "--json", "\"number,title\"", "-R", repo).Output()
	if err != nil {
		return output, err
	}

	return output, nil
}

// execEditIssueTitle uses the GitHub CLI to change the title of a given issue to the given new title
func execEditIssueTitle(repo string, number int, newTitle string) error {
	exec.Command("gh", "issue", "-R", repo)

	return nil
}

// parseIssuesToList receives a list of bytes and attempts to parse it to a list of Issues
func parseIssuesToList(issues []byte) ([]Issue, error) {
	var issueList []Issue
	if err := json.Unmarshal(issues, &issueList); err != nil {
		return nil, err
	}

	return issueList, nil
}

// editPrepend removes text prior to the first `:` found in the issueTitle,
// then prepends the prepend string along with `:`
func editPrepend(issueTitle, prepend string) string {
	titleElements := strings.SplitAfterN(issueTitle, ":", 1)

	var newTitle string
	if len(titleElements) == 2 { // replace existing prepend
		newTitle = prepend + ":" + titleElements[1]
	} else { // setting prepend where none originally present
		newTitle = prepend + ":" + titleElements[0]
	}

	return newTitle
}
