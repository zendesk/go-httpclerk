package logging

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServerLogger(t *testing.T) {
	buf, requestLogger := loadLogger()

	req, _ := http.NewRequest("PUT", "http://foo.com", new(bufio.Reader))
	res := httptest.NewRecorder()
	requestLogger.Info(res, req)

	m, _ := decodeJSONToMap(buf.String())

	if m["@timestamp"] == "" {
		t.Error("@timestamp not set correctly, expected time.RFC3339Nano, got", m["@timestamp"])
		// TODO user regex or parse to ensure correct
	}

	expectedTags := []string{"request", "foo"}
	if reflect.DeepEqual(m["@tags"], expectedTags) {
		t.Error("@tags not set correctly, expected", expectedTags, " got", m["@source"])
	}

	if m["@source"] != "fooApp" {
		t.Error("@source not set correctly, expected 'fooApp', got", m["@source"])
	}
}

func TestServerLogger_fields(t *testing.T) {
	buf, requestLogger := loadLogger()

	body := new(bufio.Reader)
	uri := "http://www.foo.com/1234.json"
	req, _ := http.NewRequest("PUT", uri, body)
	req.Header.Add("X-Foo-Header", "Bar")
	res := newWrappedRecorder()
	res.WriteHeader(http.StatusOK)
	requestLogger.Info(res, req)

	m, _ := decodeJSONToMap(buf.String())
	fields := m["@fields"].(map[string]interface{}) // Coerce again

	if fields["method"] != "PUT" {
		t.Error("@fields['method'] not set correctly, expected PUT, got", fields["method"])
	}

	if fields["status"] != "200" {
		t.Error("@fields['status'] not set correctly, expected 200, got", fields["status"])
	}

	if fields["path"] != "/1234.json" {
		t.Error("@fields['path'] not set correctly, expected /1234.json got", fields["path"])
	}

	if fields["host"] != "www.foo.com" {
		t.Error("@fields['host'] not set correctly, expected www.foo.com got", fields["host"])
	}

	headers := fields["headers"].(map[string]interface{})
	expectedHeader := []string{"Bar"}
	if reflect.DeepEqual(headers["X-Foo-Header"], expectedHeader) {
		t.Error("Header X-Foo-Header not set correctly, expected 'Bar' got", headers["X-Foo-Header"])
	}
}

func TestServerLogger_logLevelSupport_debug(t *testing.T) {
	buf, requestLogger := loadLogger()

	body := new(bufio.Reader)
	uri := "http://www.foo.com/1234.json"
	req, _ := http.NewRequest("PUT", uri, body)
	req.Header.Add("X-Foo-Header", "Bar")
	res := httptest.NewRecorder()
	requestLogger.Debug(res, req)

	m, _ := decodeJSONToMap(buf.String())
	if m["@source"] != "fooApp" {
		t.Error("Debug level not logging correctly")
	}
}

func TestServerLogger_logLevelSupport_info(t *testing.T) {
	buf, requestLogger := loadLogger()

	body := new(bufio.Reader)
	uri := "http://www.foo.com/1234.json"
	req, _ := http.NewRequest("PUT", uri, body)
	req.Header.Add("X-Foo-Header", "Bar")
	res := httptest.NewRecorder()
	requestLogger.Info(res, req)

	m, _ := decodeJSONToMap(buf.String())
	if m["@source"] != "fooApp" {
		t.Error("Info level not logging correctly")
	}
}

func TestServerLogger_logLevelSupport_warning(t *testing.T) {
	buf, requestLogger := loadLogger()

	body := new(bufio.Reader)
	uri := "http://www.foo.com/1234.json"
	req, _ := http.NewRequest("PUT", uri, body)
	req.Header.Add("X-Foo-Header", "Bar")
	res := httptest.NewRecorder()
	requestLogger.Warning(res, req)

	m, _ := decodeJSONToMap(buf.String())
	if m["@source"] != "fooApp" {
		t.Error("Warning level not logging correctly")
	}
}

func TestServerLogger_logLevelSupport_error(t *testing.T) {
	buf, requestLogger := loadLogger()

	body := new(bufio.Reader)
	uri := "http://www.foo.com/1234.json"
	req, _ := http.NewRequest("PUT", uri, body)
	req.Header.Add("X-Foo-Header", "Bar")
	res := httptest.NewRecorder()
	requestLogger.Error(res, req)

	m, _ := decodeJSONToMap(buf.String())
	if m["@source"] != "fooApp" {
		t.Error("Error level not logging correctly")
	}
}

func TestServerLogger_logLevelSupport_critical(t *testing.T) {
	buf, requestLogger := loadLogger()

	body := new(bufio.Reader)
	uri := "http://www.foo.com/1234.json"
	req, _ := http.NewRequest("PUT", uri, body)
	req.Header.Add("X-Foo-Header", "Bar")
	res := httptest.NewRecorder()
	requestLogger.Critical(res, req)

	m, _ := decodeJSONToMap(buf.String())
	if m["@source"] != "fooApp" {
		t.Error("Critical level not logging correctly")
	}
}

func TestServerLogger_responseRecorderWithoutStatusMethod(t *testing.T) {
	buf, requestLogger := loadLogger()

	body := new(bufio.Reader)
	uri := "http://www.foo.com/1234.json"
	req, _ := http.NewRequest("PUT", uri, body)
	req.Header.Add("X-Foo-Header", "Bar")
	res := httptest.NewRecorder()
	res.WriteHeader(http.StatusTeapot)
	requestLogger.Info(res, req)

	m, _ := decodeJSONToMap(buf.String())
	fields := m["@fields"].(map[string]interface{}) // Coerce again

	if fields["status"] != "" {
		t.Error("@fields['status'] not set correctly, expected 200, got", fields["status"])
	}
}

// Implement our own wrapper to set and fetch status code
type wrappedRecorder struct {
	*httptest.ResponseRecorder
}

func newWrappedRecorder() *wrappedRecorder {
	recorder := httptest.NewRecorder()
	recorder.Code = 200 // Default
	return &wrappedRecorder{recorder}
}

func (rec *wrappedRecorder) Status() int {
	return rec.Code
}

type inMemoryLogger struct {
	Buffer *bytes.Buffer
	logger *log.Logger
}

// Implement the LevelLogger interface
func (l *inMemoryLogger) Debug(info, msg string) {
	l.logger.Print(msg)
}
func (l *inMemoryLogger) Info(info, msg string) {
	l.logger.Print(msg)
}
func (l *inMemoryLogger) Warning(info, msg string) {
	l.logger.Print(msg)
}
func (l *inMemoryLogger) Error(info, msg string) {
	l.logger.Print(msg)
}
func (l *inMemoryLogger) Critical(info, msg string) {
	l.logger.Print(msg)
}

func newInMemoryLogger() *inMemoryLogger {
	buffer := new(bytes.Buffer)
	return &inMemoryLogger{Buffer: buffer, logger: log.New(buffer, "", 0)}
}

func loadLogger() (*bytes.Buffer, *HTTPLogger) {
	// Simulate a log destination we can check
	logDestination := newInMemoryLogger()

	source := "fooApp"
	tags := []string{"request", "foo"}
	formatter, _ := NewLogStashFormatter(source, tags)
	requestLogger, _ := NewHTTPLogger(logDestination, formatter)
	return logDestination.Buffer, requestLogger
}

// See: http://blog.golang.org/json-and-go ('Decoding arbitrary data' section)
func decodeJSONToMap(result string) (map[string]interface{}, error) {

	var data interface{}
	err := json.Unmarshal([]byte(result), &data)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshalling data %s", err))
	}

	return data.(map[string]interface{}), nil
}
