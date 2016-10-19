package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nezorflame/discord-wow-bot/bot"
	"github.com/nezorflame/discord-wow-bot/db"
)

var logger *log.Logger

func logInfo(v ...interface{}) {
	logger.SetPrefix("INFO  ")
	logger.Println(v...)
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
	bot.Logger = logger
	// Parse options.
	bot.DiscordToken = os.Getenv("dt")
	bot.WoWToken = os.Getenv("wt")
	bot.GoogleToken = os.Getenv("gt")
	bot.DiscordMChanID = os.Getenv("mc")
	if bot.DiscordToken == "" || bot.WoWToken == "" || bot.GoogleToken == "" || bot.DiscordMChanID == "" {
		log.Fatalln("Not enough variables to start! Abort mission! ABORT!!!")
		os.Exit(1)
	}
}

func main() {
	logInfo("Initiating db...")
	db.Init()
	logInfo("Starting bot...")
	bot.Start()
	logInfo("Bot is now running.")
	<-make(chan struct{})
}
