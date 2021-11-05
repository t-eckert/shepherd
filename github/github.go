package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Issue represents a GitHub issue
type Issue struct {
	// Unique id number for the GitHub issue
	Number int

	// Title for the GitHub issue
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

// MoveIssues iterates over the given issues and reparents them from the origin repository to the destination repository
func MoveIssues(origin, destination string, issues []Issue) error {
	for _, issue := range issues {
		if err := execTransferIssue(origin, destination, issue.Number); err != nil {
			return fmt.Errorf("could not transfer Issue %d: %s from %s to %s: %v", issue.Number, issue.Title, origin, destination, err)
		}
	}

	return nil
}

// ModifyPrepend iterates over the given issues and changes their titles to replace or add the prepended string
func ModifyPrepend(issues []Issue, repo, prepend string) error {
	for _, issue := range issues {
		newTitle := editPrepend(issue.Title, prepend)
		if err := execEditIssueTitle(repo, issue.Number, newTitle); err != nil {
			return fmt.Errorf("could not rename issue %d to %s: %v", issue.Number, newTitle, err)
		}
	}

	return nil
}

// execIssueList uses the GitHub CLI to return a list of issues
func execIssueList(repo string) ([]byte, error) {
	output, err := exec.Command("gh", "issue", "list", "--json", "number,title", "-R", repo).Output()
	if err != nil {
		return output, err
	}

	return output, nil
}

// execEditIssueTitle uses the GitHub CLI to change the title of a given issue to the given new title
func execEditIssueTitle(repo string, number int, newTitle string) error {
	if output, err := exec.Command("gh", "issue", "edit", fmt.Sprint(number), "--title", newTitle, "-R", repo).Output(); err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}

	return nil
}

// execTransferIssue uses the GitHub CLI to reparent the issue from the origin repository to the destination repository
func execTransferIssue(origin, destination string, number int) error {
	if output, err := exec.Command("gh", "issue", "transfer", fmt.Sprint(number), destination, "-R", origin).Output(); err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}

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

// editPrepend removes text prior to the first `:` found in the issueTitle, then prepends the prepend string along with `:`
func editPrepend(issueTitle, prepend string) string {
	titleElements := strings.SplitAfterN(issueTitle, ":", 2)

	var newTitle string
	if len(titleElements) == 2 { // replace existing prepend
		newTitle = prepend + ":" + titleElements[1]
	} else { // setting prepend where none originally present
		newTitle = prepend + ":" + titleElements[0]
	}

	return newTitle
}
