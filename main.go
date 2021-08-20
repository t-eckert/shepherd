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
	if prepend, ok := command.Flags["prepend"]; ok {
		// Get user confirmation before proceeding
		var shouldPrepend string
		fmt.Printf("Modify %d issues with prepend `%s:`? [y/N] ", len(issues), prepend)
		fmt.Scan(&shouldPrepend)
		if shouldPrepend != "y" || shouldPrepend != "Y" {
			log.Println("Stopping without modifying issue titles.")
			os.Exit(0)
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
	var shouldMigrate string
	fmt.Printf("Migrate %d issues from %s to %s? [y/N] ", len(issues), command.Origin, command.Destination)
	fmt.Scan(&shouldMigrate)
	if shouldMigrate != "y" || shouldMigrate != "Y" {
		log.Fatalln("Stopping without migrating issues.")
	}

	if err = github.MoveIssues(command.Origin, command.Destination, issues); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Migrated %d issues from %s to %s.", len(issues), command.Origin, command.Destination)
}
