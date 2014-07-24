package logging

import (
	"net/http"
	"strconv"
)

// Log levels to control the logging output.
const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

type LevelLogger interface {
	Debug(info, msg string)
	Info(info, msg string)
	Warning(info, msg string)
	Error(info, msg string)
	Critical(info, msg string)
}

type Formatter interface {
	Format(interface{}) (string, error)
}

type HTTPLogger struct {
	destination LevelLogger
	fmt         Formatter
}

// NewHTTPLogger constructor
func NewHTTPLogger(dest LevelLogger, formatter Formatter) (*HTTPLogger, error) {
	return &HTTPLogger{destination: dest, fmt: formatter}, nil
}

func (log *HTTPLogger) Debug(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Debug("", data)
}

func (log *HTTPLogger) Info(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Info("", data)
}

func (log *HTTPLogger) Warning(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Warning("", data)
}

func (log *HTTPLogger) Error(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Error("", data)
}

func (log *HTTPLogger) Critical(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Critical("", data)
}

// HTTP request fields that should be logged.
// If you want to log more information then add it here before setting it in
// the Log method.
type fields struct {
	Method  string              `json:"method"`
	Status  string              `json:"status"`
	Path    string              `json:"path"`
	Host    string              `json:"host"`
	Headers map[string][]string `json:"headers"`
}

func newFields(res http.ResponseWriter, req *http.Request) *fields {
	// If you need to add to the @fields key, add it here
	return &fields{
		Method:  req.Method,
		Status:  fetchStatusCode(res), // Only fetches if Status() is defined on res
		Path:    req.URL.RequestURI(),
		Headers: map[string][]string(req.Header),
		Host:    req.Host,
	}
}

// Attempts to see if the passed type implements a Status() method.
// If so, it is called and the value is returned.
// See: https://groups.google.com/forum/#!topic/golang-nuts/gz4iBqPcLt8
func fetchStatusCode(res http.ResponseWriter) string {
	var statusCode int

	type statusInterface interface {
		Status() int
	}

	statusCaller, ok := res.(statusInterface)
	if ok {
		statusCode = statusCaller.Status()
	}

	if statusCode == 0 {
		return ""
	}

	return strconv.Itoa(statusCode)
}

// Creates new fields and formats using the set formatter.
func (log *HTTPLogger) format(res http.ResponseWriter, req *http.Request) (string, error) {
	return log.fmt.Format(newFields(res, req))
}
