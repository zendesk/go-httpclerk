package httpclerk

import (
	"regexp"
	"testing"
)

func TestTextFormat(t *testing.T) {
	formatter, _ := NewTextFormatter("testApp")

	fields := &fields{
		"GET",
		"200",
		"/foo/bar",
		"dill.on.com",
		map[string][]string{"X-Foo": []string{"Gaz"}, "X-Baz": []string{"Blerg"}},
	}

	data, err := formatter.Format(fields)

	if err != nil {
		t.Error("Error formatting fields", err)
	}

	// Expecting similar to this...
	// testApp 1974-carcher.local > Method: GET Path: /foo/bar Status: 200 Host: dill.on.com Headers: map[X-Foo:[Gaz] X-Baz:[Blerg]]
	// We expect the timestamp to be prepended by the log backend.

	r, _ := regexp.Compile(`testApp`)
	if !r.MatchString(data) {
		t.Error("Ap name not formatted correctly.")
	}

	r, _ = regexp.Compile(`Host: (.)*`)
	if !r.MatchString(data) {
		t.Error("Host not formatted correctly.")
	}

	r, _ = regexp.Compile(`Method: GET`)
	if !r.MatchString(data) {
		t.Error("Method not formatted correctly.")
	}

	r, _ = regexp.Compile(`Status: 200`)
	if !r.MatchString(data) {
		t.Error("Status not formatted correctly.")
	}

	r, _ = regexp.Compile(`Path: /foo/bar`)
	if !r.MatchString(data) {
		t.Error("Path not formatted correctly.")
	}

	r, _ = regexp.Compile(`Headers: map(.)*`)
	if !r.MatchString(data) {
		t.Error("Headers not formatted correctly.")
	}

}
