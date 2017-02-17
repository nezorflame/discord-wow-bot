package main

import "github.com/golang/glog"
import "flag"

// WoWBot is a Discord WoW guild bot
var WoWBot *Bot

func main() {
	flag.Parse()
	glog.CopyStandardLogTo("INFO")
	glog.Info("Loading config...")
	LoadConfig()
	WoWBot = new(Bot)
	glog.Info("Starting...")
	WoWBot.Start()
	<-make(chan struct{})
}
