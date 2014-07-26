package httpclerk

import (
	"fmt"
	"os"
	"time"
)

type TextFormatter struct{}

func NewTextFormatter() (*TextFormatter, error) {
	return &TextFormatter{}, nil
}

func (f *fields) String() string {
	data :=
		`Method: %s
Status: %s
Path: %s
Host: %s
Headers: %s

`
	return fmt.Sprintf(data, f.Method, f.Status, f.Path, f.Host, f.Headers)
}

func (formatter *TextFormatter) Format(customFields interface{}) (string, error) {
	host, _ := os.Hostname()
	data :=
		`Time: %s
Host: %s
%s
`
	data = fmt.Sprintf(data, time.Now().Format(time.RFC3339Nano), host, customFields)

	return data, nil
}
