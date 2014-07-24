package logging

import (
	"fmt"
	golog "github.com/op/go-logging"
	stdlog "log"
	"net/http"
	"os"
	"strconv"
)

const (
	BackendStdOut = iota
	BackendSysLog
	BackendMemory
)

type Formatter interface {
	Format(interface{}) (string, error)
}

type HTTPLogger struct {
	name        string
	Backends    []int
	formatter   Formatter
	destination *golog.Logger
}

// NewHTTPLogger constructor
func NewHTTPLogger(name string, backends []int, formatter Formatter) (*HTTPLogger, error) {
	destination := golog.MustGetLogger(name)

	// Customize the output format
	// golog.SetFormatter(golog.MustStringFormatter("▶ %{level:.1s} 0x%{id:x} %{message}"))

	requiredBackends := make([]golog.Backend, 0)
	var logBackend golog.Backend

	for _, backend := range backends {
		switch backend {
		case BackendStdOut:
			fmt.Println("stdout")
			logBackend = golog.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
			requiredBackends = append(requiredBackends, logBackend)
		case BackendSysLog:
			fmt.Println("syslog")
			logBackend, err := golog.NewSyslogBackend(name)
			if err != nil {
				stdlog.Fatal("Could not setup syslog backend.", err)
			}
			requiredBackends = append(requiredBackends, logBackend)
		case BackendMemory:
			fmt.Println("mem")
			logBackend = golog.NewMemoryBackend(1024)

			requiredBackends = append(requiredBackends, logBackend)
		}

	}

	fmt.Println("len backends", len(requiredBackends))
	if len(requiredBackends) == 0 {
		stdlog.Fatal("Please supply at least one backend!")
	}

	// Set the backend to whatever backends were specified
	golog.SetBackend(requiredBackends...)

	return &HTTPLogger{name: name, Backends: backends, formatter: formatter, destination: destination}, nil

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
	return log.formatter.Format(newFields(res, req))
}
