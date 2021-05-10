package main

import "testing"

func TestFilter(t *testing.T) {
	shouldBeEmpty := [...]string{
		`<a>`,
		`https://www.alexrecker.com`,
		`src="images/bob.jpg"`,
		`?`,
		`.`,
		`-`,
	}

	for _, token := range shouldBeEmpty {
		actual := Filter(token)
		if actual != "" {
			t.Errorf("expected %s to be empty, got %s", token, actual)
		}
	}

}
