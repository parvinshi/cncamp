package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/parvinshi/cncamp/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting HTTP server...")
	metrics.Register()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	addrPort := os.Getenv("ADDR_PORT")
	glog.Infof("Env variable ADDR_PORT: %v", addrPort)
	if addrPort == "" {
		addrPort = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/foo", fooHandler)
	mux.HandleFunc("/bar", barHandler)
	mux.HandleFunc("/foobar", foobarHandler)
	mux.HandleFunc("/latency", latencyHandler)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", healthzHandler)
	srv := &http.Server{
		Addr:    ":" + addrPort,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("HTTP ListenAndServe Error: %s\n", err)
		}
	}()
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
	glog.V(4).Info("Entering index handler")
	headers := make(http.Header)
	for k, v := range r.Header {
		headers[strings.ToLower(k)] = v

		w.Header().Add(k, r.Header.Get(k))
	}

	glog.Infof("client ip:port -> %v", r.RemoteAddr)
	glog.Infof("headers:%+v", headers)
}

func latencyHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("Entering latency handler")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	randInt := rand.Intn(2000)
	time.Sleep(time.Millisecond * time.Duration(randInt))
	w.Write([]byte(fmt.Sprintf("<h1>latency %d<h1>", randInt)))
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("Entering foo handler")
	headers := make(http.Header)
	for k, v := range r.Header {
		headers[strings.ToLower(k)] = v

		w.Header().Add(k, r.Header.Get(k))
	}

	glog.Infof("client ip:port -> %v", r.RemoteAddr)
	glog.Infof("headers:%+v", headers)

	w.Write([]byte("foo"))
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("Entering bar handler")
	headers := make(http.Header)
	for k, v := range r.Header {
		headers[strings.ToLower(k)] = v

		w.Header().Add(k, r.Header.Get(k))
	}

	glog.Infof("client ip:port -> %v", r.RemoteAddr)
	glog.Infof("headers:%+v", headers)

	w.Write([]byte("bar"))
}

func foobarHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("Entering foobar handler")
	headers := make(http.Header)
	for k, v := range r.Header {
		headers[strings.ToLower(k)] = v

		w.Header().Add(k, r.Header.Get(k))
	}

	glog.Infof("client ip:port -> %v", r.RemoteAddr)
	glog.Infof("headers:%+v", headers)

	w.Write([]byte("foobar"))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("Entering healthz handler")
	glog.Infof("client request uri: %v", r.RequestURI)

	w.WriteHeader(http.StatusOK)
}
