package httpclerk

import (
	"regexp"
	"testing"
)

func TestTextFormat(t *testing.T) {
	formatter, _ := NewTextFormatter()

	fields := &fields{
		"GET",
		"200",
		"/foo/bar",
		"dill.on.com",
		map[string][]string{"X-Foo": []string{"Gaz"}, "X-Baz": []string{"Blerg"}},
	}

	data, err := formatter.Format(fields)

	// Expecting similar to this...
	//
	// Time: 2014-07-26T21:37:14.149045474+01:00
	// Host: 1974-carcher.local
	// Method: GET
	// Status: 200
	// Path: /foo/bar
	// Host: dill.on.com
	// Headers: map[X-Foo:[Gaz] X-Baz:[Blerg]]

	if err != nil {
		t.Error("Error formatting fields", err)
	}

	r, _ := regexp.Compile(`Host: (.)*`)
	if !r.MatchString(data) {
		t.Error("Host not formatted correctly.")
	}

	r, _ = regexp.Compile(`Time: ((.)*T(.)*)`)
	if !r.MatchString(data) {
		t.Error("Time not formatted correctly.")
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
