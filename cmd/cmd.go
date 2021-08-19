package cmd

import "fmt"

// Command represents arguments and flags passed into the application
type Command struct {
	// GitHub repository issues will be migrated from
	Origin string

	// GitHub repository issues will be migrated to
	Destination string

	// Optional flags passed in
	Flags map[string]string
}

var validFlags = map[string]struct {
	description   string
	alias         string
	requiresValue bool
}{
	"--modify-prepend": {"Change the string prepended to each GitHub issue (e.g. helm:). Takes the new string to prepend as an argument.", "", true},
	"-p":               {"", "--modify-prepend", true},
}

// NewCommand parses arguments passed into the application and returns a Command object
func NewCommand(args []string) (Command, error) {
	if err := validateArgs(args); err != nil {
		return Command{}, fmt.Errorf("could not parse command line arguments: %v", err)
	}

	if err := validateFlags(args); err != nil {
		return Command{}, fmt.Errorf("could not parse flags: %v", err)
	}

	origin, destination := args[1], args[2]

	flags := parseFlags(args)

	return Command{origin, destination, flags}, nil
}

func validateArgs(args []string) error {
	if len(args) == 1 {
		return fmt.Errorf("missing argument for origin repository")
	} else if len(args) == 2 {
		return fmt.Errorf("missing argument for destination repository")
	}

	return nil
}

func validateFlags(args []string) error {
	for _, arg := range args {
		if _, ok := validFlags[arg]; arg[0] == '-' && !ok {
			return fmt.Errorf("unknown flag %s passed", arg)
		}
	}

	return nil
}

func parseFlags(args []string) map[string]string {
	flags := make(map[string]string)
	for i, arg := range args {
		if arg[0] != '-' { // not a flag
			continue
		}

		flag, value := resolveFlag(arg), args[i+1]
		flags[flag] = value
	}

	return flags
}

func resolveFlag(flag string) string {
	if alias := validFlags[flag].alias; alias != "" {
		flag = alias
	}

	return flag[2:]
}
