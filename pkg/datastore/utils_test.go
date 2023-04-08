package datastore

import "testing"

func TestParseDbUrl(t *testing.T) {

	testStrings := []struct {
		input    string
		expected string
	}{{
		"sqlite3://test.db",
		"sqlite",
	}, {
		"mysql://user:password@tcp(mariadb:3306)",
		"mysql",
	}}

	for _, test := range testStrings {
		t.Run(test.input, func(t *testing.T) {
			d, err := ParseDBUrl(test.input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if d.Name() != test.expected {
				t.Errorf("expected %s, got %s", test.expected, d.Name())
			}

		})
	}

}
