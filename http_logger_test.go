package httpclerk

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	golog "github.com/op/go-logging"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServerLogger_logLevelSupport_debug(t *testing.T) {
	_, logger := loadLogger()
	res, req := createRequestAndResponse()

	logger.Debug(res, req)
}

func TestServerLogger_logLevelSupport_info(t *testing.T) {
	_, logger := loadLogger()
	res, req := createRequestAndResponse()

	logger.Info(res, req)
}

func TestServerLogger_logLevelSupport_warning(t *testing.T) {
	_, logger := loadLogger()
	res, req := createRequestAndResponse()

	logger.Warning(res, req)
}

func TestServerLogger_logLevelSupport_error(t *testing.T) {
	_, logger := loadLogger()
	res, req := createRequestAndResponse()

	logger.Error(res, req)
}

func TestServerLogger_logLevelSupport_critical(t *testing.T) {
	_, logger := loadLogger()
	res, req := createRequestAndResponse()

	logger.Critical(res, req)
}

func TestLogger_fields(t *testing.T) {
	memBackend, logger := loadLogger()
	res, req := createRequestAndResponse()
	res = newWrappedRecorder()
	res.WriteHeader(http.StatusOK)

	logger.Info(res, req)

	lastWrite := memBackend.Head().Record.Message()
	m, _ := decodeJSONToMap(lastWrite)
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

func TestServerLogger_responseRecorderWithoutStatusMethod(t *testing.T) {
	memBackend, logger := loadLogger()
	res, req := createRequestAndResponse()
	res = httptest.NewRecorder()
	res.WriteHeader(http.StatusTeapot)

	logger.Info(res, req)
	lastWrite := memBackend.Head().Record.Message()

	m, _ := decodeJSONToMap(lastWrite)
	fields := m["@fields"].(map[string]interface{}) // Coerce again

	if fields["status"] != "" {
		t.Error("@fields['status'] not set correctly, expected 200, got", fields["status"])
	}
}

// *************************************
// Helper functions
// *************************************

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func loadLogger() (*golog.MemoryBackend, *HTTPLogger) {
	log := golog.MustGetLogger("test")
	formatter, _ := NewLogStashFormatter("fooApp", []string{"blimp", "foo"})

	// Setup one stdout and one syslog backend.
	memBackend := golog.NewMemoryBackend(1024)

	// Combine them both into one logging backend and set the log level
	golog.SetBackend(memBackend)

	logger, _ := NewHTTPLogger("foo", log, formatter)
	return memBackend, logger
}

func createRequestAndResponse() (http.ResponseWriter, *http.Request) {
	body := new(bufio.Reader)
	uri := "http://www.foo.com/1234.json"
	req, _ := http.NewRequest("PUT", uri, body)
	req.Header.Add("X-Foo-Header", "Bar")
	res := httptest.NewRecorder()
	return res, req
}

// // Implement our own wrapper to set and fetch status code
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

// See: http://blog.golang.org/json-and-go ('Decoding arbitrary data' section)
func decodeJSONToMap(result string) (map[string]interface{}, error) {

	var data interface{}
	err := json.Unmarshal([]byte(result), &data)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshalling data %s", err))
	}

	return data.(map[string]interface{}), nil
}
