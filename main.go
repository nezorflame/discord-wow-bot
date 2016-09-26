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

func determineListenAddress() (string, error) {
  port := os.Getenv("PORT")
  if port == "" {
    return "", fmt.Errorf("$PORT not set")
  }
  return ":" + port, nil
}

func watcherHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Starting guild watcher...")
    go bot.RunGuildWatcher()
}

func aliveHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hi there! I'm alive! :D")
}

func httpStart(addr string) {
    http.HandleFunc("/", aliveHandler)
    http.HandleFunc("/startwatcher", watcherHandler)
    if err := http.ListenAndServe(addr, nil); err != nil {
        panic(err)
    }
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
    addr, err := determineListenAddress()
    if err != nil {
        log.Fatal(err)
    }
    bot.Start()
    logInfo("Starting handler...")
	httpStart(addr)
	logInfo("Bot is now running.")
}
