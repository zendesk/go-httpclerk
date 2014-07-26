package httpclerk

import (
	"reflect"
	"testing"
)

type testingFields struct {
	Hot   string
	Chip  string
	Block int
	Blip  []int
}

func TestLogstashFormat(t *testing.T) {
	formatter, _ := NewLogStashFormatter("fooApp", []string{"blimp", "foo"})

	fields := &testingFields{"hi", "there", 101, []int{1, 2, 3, 4}}
	data, err := formatter.Format(fields)

	if err != nil {
		t.Error("Error formatting fields", err)
	}

	m, _ := decodeJSONToMap(data)

	if m["@timestamp"] == "" {
		t.Error("@timestamp not set correctly, expected time.RFC3339Nano, got", m["@timestamp"])
	}

	expectedTags := []string{"request", "foo"}
	if reflect.DeepEqual(m["@tags"], expectedTags) {
		t.Error("@tags not set correctly, expected", expectedTags, " got", m["@source"])
	}

	if m["@source"] != "fooApp" {
		t.Error("@source not set correctly, expected 'fooApp', got", m["@source"])
	}
}
