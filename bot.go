package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/golang/glog"
)

// Bot struct for the Discord bot
type Bot struct {
	ID      string
	Session *discordgo.Session
}

// Start launches the connection to the bot
func (b *Bot) Start() {
	var (
		err error
		u   *discordgo.User
	)

	glog.Info("Logging in...")
	if b.Session, err = discordgo.New(o.DiscordToken); err != nil {
		glog.Fatalf("Unable to connect to Discord: %s", err)
	}

	glog.Info("Using bot account token...")
	if u, err = b.Session.User("@me"); err != nil {
		glog.Fatalf("Unable to get @me: %s", err)
	} else {
		b.ID = u.ID
		glog.Infof("Got BotID = %s", b.ID)
	}

	glog.Info("Adding handlers...")
	b.Session.AddHandler(b.messageCreate)

	glog.Info("Opening session...")
	if err = b.Session.Open(); err != nil {
		glog.Fatalf("Unable to open the session: %s", err)
	}

	glog.Info("Starting guild watcher...")
	// legendaries = make(map[string][]*Item)
	// go guildWatcher(session)

	glog.Info("Bot started")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (b *Bot) messageCreate(s *discordgo.Session, mes *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if mes.Author.ID == b.ID {
		return
	}
}
