package main

import "github.com/golang/glog"

// WoWBot is a Discord WoW guild bot
var WoWBot *Bot

func main() {
	glog.CopyStandardLogTo("INFO")
	glog.Info("Loading config...")
	LoadConfig()
	WoWBot = new(Bot)
	glog.Info("Starting...")
	WoWBot.Start()
	<-make(chan struct{})
}
