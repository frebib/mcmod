package minecraft

import (
	"reflect"
	"testing"
)

func TestVersions(t *testing.T) {
	var cases = []struct {
		input    string
		normal   string
		expected Version
	}{
		{input: "1.15.2", expected: Version{Release: Release, Major: 1, Minor: 15, Patch: 2, Build: -1}},
		{input: "1.14", expected: Version{Release: Release, Major: 1, Minor: 14, Patch: -1, Build: -1}},
		{input: "1.7.10", expected: Version{Release: Release, Major: 1, Minor: 7, Patch: 10, Build: -1}},
		{input: "1.0_01", expected: Version{Release: Release, Major: 1, Minor: 0, Patch: -1, Build: 1}},
		{input: "1.0", expected: Version{Release: Release, Major: 1, Minor: 0, Patch: -1, Build: -1}},
		{input: "Beta 1.5_01", expected: Version{Release: Beta, Major: 1, Minor: 5, Patch: -1, Build: 1}},
		{
			input:    "Alpha v1.0.17_04",
			normal:   "Alpha 1.0.17_04",
			expected: Version{Release: Alpha, Major: 1, Minor: 0, Patch: 17, Build: 4},
		},
		{
			input:    "Alpha v1.1.2_01",
			normal:   "Alpha 1.1.2_01",
			expected: Version{Release: Alpha, Major: 1, Minor: 1, Patch: 2, Build: 1},
		},
	}

	for _, c := range cases {
		parsed, err := Parse(c.input)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(parsed, &c.expected) {
			t.Errorf("unexpected parsed output: %s\nexpected: %#v\ngot:     %#v", c.input, &c.expected, parsed)
		}
		if c.normal == "" {
			c.normal = c.input
		}
		if c.normal != parsed.String() {
			t.Errorf("input not equal to parsed string output: %s != %s", c.input, parsed.String())
		}
	}
}
