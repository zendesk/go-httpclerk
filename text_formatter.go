package httpclerk

import (
	"fmt"
	"os"
)

type TextFormatter struct {
	AppName string
}

func NewTextFormatter() (*TextFormatter, error) {
	return &TextFormatter{"testApp"}, nil
}

func (f *fields) String() string {
	data := `Method: %s Path: %s Status: %s Host: %s Headers: %s`
	return fmt.Sprintf(data, f.Method, f.Path, f.Status, f.Host, f.Headers)
}

func (f *TextFormatter) Format(customFields interface{}) (string, error) {
	host, _ := os.Hostname()
	data := `%s %s > %s`
	data = fmt.Sprintf(data, f.AppName, host, customFields)

	return data, nil
}
