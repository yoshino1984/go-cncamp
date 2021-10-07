package interceptor

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

func HandleLogInterceptor(h http.HandlerFunc) http.HandlerFunc {
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
