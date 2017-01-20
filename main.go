package main

import (
	"fmt"
	"log"
	"os"
)

var (
	// Logger for logging
	Logger *log.Logger
	// DiscordToken bot token
	DiscordToken string
	// WoWToken API token
	WoWToken string
	// GoogleToken API token
	GoogleToken string
	// DiscordMChanID main guild channel ID
	DiscordMChanID string
	// GuildRosterMID guild roster message ID
	GuildRosterMID string
	// GuildName guild name
	GuildName string
	// GuildRealm guild realm
	GuildRealm string
)

func logDebug(v ...interface{}) {
	Logger.SetPrefix("DEBUG ")
	Logger.Println(v...)
}

func logInfo(v ...interface{}) {
	Logger.SetPrefix("INFO  ")
	Logger.Println(v...)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func logOnErr(err error) {
	if err != nil {
		logDebug(err)
	}
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func init() {
	// Create initials.
	Logger = log.New(os.Stderr, "  ", log.Ldate|log.Ltime)
	// Parse options.
	DiscordToken = os.Getenv("dt")
	WoWToken = os.Getenv("wt")
	GoogleToken = os.Getenv("gt")
	DiscordMChanID = os.Getenv("mc")
	GuildRosterMID = os.Getenv("gr")
	GuildName = os.Getenv("gn")
	GuildRealm = os.Getenv("re")
	if DiscordToken == "" || WoWToken == "" || GoogleToken == "" || DiscordMChanID == "" ||
		GuildRosterMID == "" || GuildName == "" || GuildRealm == "" {
		log.Fatalln("Not enough variables to start! Abort mission! ABORT!!!")
		os.Exit(1)
	}
}

func main() {
	logInfo("Initiating ..")
	Init()
	defer Close()
	go Watcher()
	logInfo("Starting ..")
	Start()
	logInfo("Bot is now running.")
	<-make(chan struct{})
}
