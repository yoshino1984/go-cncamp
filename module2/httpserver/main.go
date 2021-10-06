package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	//http.HandleFunc("/", defaultHandler)
	//http.HandleFunc("/healthz", healthzHandler)
	addHandleFunc("/", defaultHandler)
	addHandleFunc("/healthz", healthzHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		glog.V(4).Info("start listen")
	}
}

func addHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handleLogInterceptor(handler))
}

type doneWriter struct {
	http.ResponseWriter
	done   bool
	status int
}

func (w *doneWriter) WriteHeader(statusCode int) {
	w.done = true
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
func (w *doneWriter) Write(b []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(b)
}

func handleLogInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dw := &doneWriter{ResponseWriter: w}
		h(dw, r)
		if dw.done {
			ip, _ := GetIP(r)
			code := 200
			if dw.status != 0 {
				code = dw.status
			}
			msg := "IP:" + ip + ", STATUS:" + strconv.Itoa(code)
			os.Stdout.WriteString("\nlog: " + msg)
		}
	}
}

func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}

func defaultHandler(writer http.ResponseWriter, request *http.Request) {
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

func healthzHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	//io.WriteString(writer, "healthz")
}
