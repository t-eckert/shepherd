package cmd

import (
	"fmt"
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
			Command{"organization/origin", "organization/destination", map[string]string{"modify-prepend": "prepend"}},
			nil,
		},
		{
			[]string{"main.go", "organization/origin", "organization/destination", "--modify-prepend", "prepend"},
			Command{"organization/origin", "organization/destination", map[string]string{"modify-prepend": "prepend"}},
			nil,
		},
	}

	for _, tt := range params {
		actualCommand, actualError := NewCommand(tt.given)

		if actualCommand.Origin != tt.expectedCommand.Origin || actualCommand.Destination != tt.expectedCommand.Destination {
			fmt.Printf("Given: %v\nExpected: %v\nReceived: %v\n", tt.given, tt.expectedCommand, actualCommand)
			t.Fail()
		}

		for expectedKey, expectedValue := range tt.expectedCommand.Flags {
			if actualValue, ok := actualCommand.Flags[expectedKey]; !ok || actualValue != expectedValue {
				fmt.Printf("Given: %v\nExpected: %v\nReceived: %v\n", tt.given, expectedValue, actualValue)
				t.Fail()
			}
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

func TestValidateFlags(t *testing.T) {
	params := []struct {
		given    []string
		expected error
	}{
		{[]string{"main.go"}, nil},
		{[]string{"main.go", "organization/origin", "organization/destination", "-p"}, fmt.Errorf("missing value for flag -p")},
		{[]string{"main.go", "organization/origin", "organization/destination", "-x"}, fmt.Errorf("unknown flag -x passed")},
		{[]string{"main.go", "organization/origin", "organization/destination", "-p", "prepend"}, nil},
		{[]string{"main.go", "organization/origin", "organization/destination", "--modify-prepend", "prepend"}, nil},
	}

	for _, tt := range params {
		actual := validateFlags(tt.given)

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
		{[]string{"main.go", "organization/origin", "organization/destination", "-p", "prepend"}, map[string]string{"modify-prepend": "prepend"}},
		{[]string{"main.go", "organization/origin", "organization/destination", "--modify-prepend", "prepend"}, map[string]string{"modify-prepend": "prepend"}},
	}

	for _, tt := range params {
		actual := parseFlags(tt.given)

		for expectedKey, expectedValue := range tt.expected {
			if actualValue, ok := actual[expectedKey]; !ok || actualValue != expectedValue {
				fmt.Printf("Given: %v\nExpected: %v\nReceived: %v\n", tt.given, expectedValue, actualValue)
				t.Fail()
			}
		}
	}
}
