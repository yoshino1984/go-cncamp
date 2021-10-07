package main

import (
	"flag"
	"github.com/golang/glog"
	"log"
	"net/http"
	"yoshino.com/cncamp/module2/httpserver/handler"
	"yoshino.com/cncamp/module2/httpserver/interceptor"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	addHandleFunc("/", handler.DefaultHandler)
	addHandleFunc("/healthz", handler.HealthzHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		glog.V(4).Info("start listen")
	}
}

func addHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, interceptor.HandleLogInterceptor(handler))
}
