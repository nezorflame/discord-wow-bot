package main

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

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
	logger = log.New(os.Stderr, "  ", log.Ldate|log.Ltime)
	Logger = logger
	// Parse options.
	DiscordToken = os.Getenv("dt")
	WoWToken = os.Getenv("wt")
	GoogleToken = os.Getenv("gt")
	DiscordMChanID = os.Getenv("mc")
	if DiscordToken == "" || WoWToken == "" || GoogleToken == "" || DiscordMChanID == "" {
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
