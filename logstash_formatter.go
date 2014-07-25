package scribe

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type LogStashFormatter struct {
	Source string
	Tags   []string
}

func NewLogStashFormatter(source string, tags []string) (*LogStashFormatter, error) {
	return &LogStashFormatter{Source: source, Tags: tags}, nil
}

type LogStashJSON struct {
	Source    string      `json:"@source"`
	Fields    interface{} `json:"@fields"`
	Tags      []string    `json:"@tags"`
	Timestamp string      `json:"@timestamp"`
}

func (formatter *LogStashFormatter) Format(customFields interface{}) (string, error) {
	stash := &LogStashJSON{
		Source:    formatter.Source,
		Fields:    customFields,
		Tags:      formatter.Tags,
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}

	data, err := json.Marshal(stash)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error marshalling JSON: %s", err))
	}

	return string(data), nil
}
