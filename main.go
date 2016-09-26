package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/nezorflame/discord-wow-bot/bot"
)

var logger *log.Logger

func logInfo(v ...interface{}) {
	logger.SetPrefix("INFO  ")
	logger.Println(v...)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func httpStart() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
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
    logInfo("Starting bot...")
    bot.Start()
    logInfo("Starting handler...")
	httpStart()
	logInfo("Bot is now running.\nPress CTRL-C to exit...")
	<-make(chan struct{})
}
