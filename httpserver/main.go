package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting HTTP server...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	srv := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	go srv.ListenAndServe()
	glog.Info("Server started")

	<-sig
	glog.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		glog.Fatalf("Server Shutdown Faliled: %+v", err)
	}

	glog.Info("Server Shutdown successfully")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("enter index handler")
	headers := make(http.Header)
	for k, v := range r.Header {
		headers[strings.ToLower(k)] = v

		w.Header().Add(k, r.Header.Get(k))
	}

	glog.Infof("client ip:port -> %v", r.RemoteAddr)
	glog.Infof("headers:%+v", headers)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("enter healthz handler")
	glog.Infof("client request uri: %v", r.RequestURI)

	w.WriteHeader(http.StatusOK)
}
