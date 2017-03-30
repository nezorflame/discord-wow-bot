package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/golang/glog"
)

// WoWBot is a Discord WoW guild bot
var (
	WoWBot *Bot
)

func main() {
	flag.Parse()
	glog.CopyStandardLogTo("INFO")
	glog.Info("Loading config...")
	LoadConfig()
	WoWBot = new(Bot)
	glog.Info("Starting...")
	WoWBot.Start()
	log.Println(http.ListenAndServe("localhost:6060", nil))
	<-make(chan struct{})
}
