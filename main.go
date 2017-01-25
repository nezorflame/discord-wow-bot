package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {
	flag.Parse()
	glog.CopyStandardLogTo("INFO")
	glog.Info("Loading config...")
	LoadConfig()
	glog.Info("Initiating DB connection...")
	InitDB()
	defer CloseDB()
	go DBWatcher()
	glog.Info("Starting...")
	Start()
	<-make(chan struct{})
}
