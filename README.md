go-httpclerk
==============

## Overview

A simple HTTP request/response logger for Go supporting multiple formatters.

## Rationale

We needed a way to log HTTP requests at Zendesk to different log backends (stdout, syslog etc.) with multiple ways to format them (including logstash). So we created this project to help us. 

## Usage

You'll need to create some sort of logger that conforms to the `LogDestination` interface in this package. The [go-logger](https://github.com/op/go-logging) package is recommended.

### Simple example:

```
package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/zendesk/go-httpclerk"
	stdlog "log"
	"net/http"
	"os"
)

var log = logging.MustGetLogger("myApp")

func main() {
	// Setup a go-logging logger
	stdoutBackend := logging.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
	logging.SetBackend(stdoutBackend) // See go-logging docs for multiple backends
	logging.SetLevel(logging.DEBUG, "myApp")

	// Boot web server and listen on 8080
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	formatter, _ := httpclerk.NewTextFormatter("myHandler")
	clerk, err := httpclerk.NewHTTPLogger("myHandler", log, formatter)
	if err != nil {
		log.Fatal("HTTP logger could not be created", err)
	}
	defer clerk.Info(w, r)

	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
```

This will produce logs like so: 

```
2014/07/27 07:43:56 http_logger.go:39: myHandler 1974-carcher.local > Method: GET Path: /ciaran Status:  Host: localhost:8080 Headers: map[User-Agent:[curl/7.30.0] Accept:[*/*]]
```

### Status Code

You'll notice that the `Status` is blank. This is becuase there is no simple way to get the response status in a HTTP handler without wrapping the `ResponseWriter` type. You can see an example of this [here](https://gist.github.com/ciaranarcher/abccf50cb37645ca27fa). If you do this, and use this type instead of the standard `ResponseWriter` then the `go-httpclerk` package can fetch the status code and include it in logging:

```
2014/07/27 07:43:56 http_logger.go:39: myHandler 1974-carcher.local > Method: GET Path: /ciaran Status: 200 Host: localhost:8080 Headers: map[User-Agent:[curl/7.30.0] Accept:[*/*]]
```

### Other Formatters

Included in the package is a `TextFormatter` (examples above use this) and a `LogStashFormatter` for JSON logging

```
formatter, _ := NewLogStashFormatter("fooApp", []string{"blimp", "foo"})
```

Other loggers can be used in place if they implement the the following interface:

```
type Formatter interface {
	Format(interface{}) (string, error)
}
```

## Contributing

Create a Pull Request with your changes, ping someone and we'll look at getting it merged. 

