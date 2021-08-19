package cmd

import (
	"fmt"
	"testing"
)

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
