package handler

import (
	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
	"os"
)

func DefaultHandler(writer http.ResponseWriter, request *http.Request) {
	glog.V(5).Info("entering root handler")
	processHeader(writer, request)
	user := request.URL.Query().Get("user")
	if user != "" {
		io.WriteString(writer, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(writer, "hello [stranger]\n")
	}
	io.WriteString(writer, "===================Details of the http request header:============\n")
	for k, v := range request.Header {
		io.WriteString(writer, fmt.Sprintf("%s=%s\n", k, v))
	}
}

func HealthzHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	//io.WriteString(writer, "healthz")
}

func processHeader(writer http.ResponseWriter, request *http.Request) {
	for k, v := range request.Header {
		for i := range v {
			if i == 0 {
				writer.Header().Set(k, v[i])
			} else {
				writer.Header().Add(k, v[i])
			}
		}
	}
	writer.Header().Set("X-Hello", "yoshino")
	writer.Header().Set("Version", os.Getenv("VERSION"))
	writer.WriteHeader(http.StatusOK)
}
