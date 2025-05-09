package main

import (
    "net/http"
    "os"

    log "github.com/sirupsen/logrus"
)

func init() {
    log.SetFormatter(&log.JSONFormatter{})
    log.SetLevel(log.InfoLevel)
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.WithFields(log.Fields{
        "method": r.Method,
        "ip":     r.RemoteAddr,
    }).Info("Incoming request")
    w.Write([]byte("Hello from Go with logrus!"))
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.WithField("port", port).Info("Starting server")

    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        log.WithError(err).Fatal("Server failed to start")
    }
}