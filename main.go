package main

import (
	"fmt"
	"log"
	"os"

	"github.com/t-eckert/shepherd/cmd"
	"github.com/t-eckert/shepherd/github"
)

func main() {
	args := os.Args

	command, err := cmd.NewCommand(args)
	if err != nil {
		log.Fatalln(err)
	}

	issues, err := github.FetchIssues(command.Origin)
	if err != nil {
		log.Fatalln(err)
	}

	// Modify the prepend (e.g. helm:) if requested
	if prepend, ok := command.Flags["modify-prepend"]; ok {
		// Get user confirmation before proceeding
		shouldPrepend := askToContinue(fmt.Sprintf("Modify %d issues with prepend `%s:`?", len(issues), prepend))
		if !shouldPrepend {
			safeAbort("Stopping without modifying issue titles.")
		}

		if err := github.ModifyPrepend(issues, command.Origin, prepend); err != nil {
			log.Fatalln(err)
		}

		// Refetch the updated issues
		issues, err = github.FetchIssues(command.Origin)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Modified %d issues", len(issues))
	}

	// Get user confirmation before proceeding
	shouldMigrate := askToContinue(fmt.Sprintf("Migrate %d issues from %s to %s? [y/N] ", len(issues), command.Origin, command.Destination))
	if !shouldMigrate {
		safeAbort("Stopping without migrating issues.")
	}

	if err = github.MoveIssues(command.Origin, command.Destination, issues); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Migrated %d issues from %s to %s.", len(issues), command.Origin, command.Destination)
}

// askToContinue prompts the user with the given message,
// if user responds with "y", "Y", "yes", or "Yes", returns true
func askToContinue(message string) bool {
	affirmativeResponses := []string{"y", "Y", "yes", "Yes"}
	var response string

	fmt.Printf(message + " [y/N]")
	fmt.Scanln(&response)

	for _, affirmativeResponse := range affirmativeResponses {
		if response == affirmativeResponse {
			return true
		}
	}

	return false
}

// safeAbort prints the given message then exits the program
func safeAbort(message string) {
	log.Println(message)
	os.Exit(0)
}
