package cmd

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewCommand(t *testing.T) {
	params := []struct {
		given           []string
		expectedCommand Command
		expectedError   error
	}{
		{
			[]string{"main.go"},
			Command{},
			fmt.Errorf("missing argument for origin repository"),
		},
		{
			[]string{"main.go", "organization/origin"},
			Command{},
			fmt.Errorf("missing argument for destination repository"),
		},
		{
			[]string{"main.go", "organization/origin", "organization/destination"},
			Command{"organization/origin", "organization/destination", map[string]string{}},
			nil,
		},
		{
			[]string{"main.go", "organization/origin", "organization/destination", "-p"},
			Command{},
			fmt.Errorf("missing value for flag -p"),
		},
		{
			[]string{"main.go", "organization/origin", "organization/destination", "-p", "prepend"},
			Command{"organization/origin", "organization/destination", map[string]string{"prepend": "prepend"}},
			nil,
		},
		{
			[]string{"main.go", "organization/origin", "organization/destination", "--modify-prepend", "prepend"},
			Command{"organization/origin", "organization/destination", map[string]string{"prepend": "prepend"}},
			nil,
		},
	}

	for _, tt := range params {
		actualCommand, actualError := NewCommand(tt.given)

		if actualCommand.Origin != tt.expectedCommand.Origin ||
			actualCommand.Destination != tt.expectedCommand.Destination ||
			reflect.DeepEqual(actualCommand.Flags, tt.expectedCommand.Flags) {
			fmt.Printf("Given: %v\nExpected: %v\nReceived: %v\n", tt.given, tt.expectedCommand, actualCommand)
			t.Fail()
		}

		if actualError != nil && tt.expectedError == nil {
			fmt.Printf("No error expected. Received error %s\n", actualError)
			t.Fail()
		} else if actualError == nil && tt.expectedError != nil {
			fmt.Printf("Error expected %s. No error received.\n", tt.expectedError)
			t.Fail()
		}
	}
}

func TestValidateArgs(t *testing.T) {
	params := []struct {
		given    []string
		expected error
	}{
		{[]string{"main.go"}, fmt.Errorf("missing argument for origin repository")},
		{[]string{"main.go", "organization/origin"}, fmt.Errorf("missing argument for destination repository")},
		{[]string{"main.go", "organization/origin", "organization/destination"}, nil},
	}

	for _, tt := range params {
		actual := validateArgs(tt.given)

		if actual != nil && tt.expected == nil {
			fmt.Printf("No error expected. Received error %s", actual)
			t.Fail()
		} else if actual == nil && tt.expected != nil {
			fmt.Printf("Error expected %s. No error received.", tt.expected)
			t.Fail()
		}
	}
}

func TestParseFlags(t *testing.T) {
	params := []struct {
		given    []string
		expected map[string]string
	}{
		{[]string{"main.go"}, map[string]string{}},
		{[]string{"main.go", "organization/origin", "organization/destination"}, map[string]string{}},
		{[]string{"main.go", "organization/origin", "organization/destination", "-p"}, map[string]string{}},
		{[]string{"main.go", "organization/origin", "organization/destination", "-p", "prepend"}, map[string]string{"prepend": "prepend"}},
		{[]string{"main.go", "organization/origin", "organization/destination", "--modify-prepend", "prepend"}, map[string]string{"prepend": "prepend"}},
	}

	for _, tt := range params {
		actual := parseFlags(tt.given)

		if reflect.DeepEqual(actual, tt.expected) {
			fmt.Printf("Given: %s Expected: %s Received: %s", tt.given, tt.expected, actual)
			t.Fail()
		}
	}
}

