package main

import (
	"strings"
	"time"

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
	b.Session.AddHandler(b.parseMessage)

	glog.Info("Opening session...")
	if err = b.Session.Open(); err != nil {
		glog.Fatalf("Unable to open the session: %s", err)
	}

	glog.Info("Starting guild watcher...")
	// legendaries = make(map[string][]*Item)
	// go guildWatcher(session)

	glog.Info("Bot started")
}

// SendMessage sends the message to the selected channel
func (b *Bot) SendMessage(chID string, message string) (err error) {
	glog.Info("SENDING MESSAGE:", message)
	retryOnBadGateway(func() error {
		return sendFormattedMessage(b.Session, chID, message)
	})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (b *Bot) parseMessage(s *discordgo.Session, mes *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself or without exclamation mark
	if mes.Author.ID == b.ID || !strings.HasPrefix(mes.Content, "!") {
		return
	}
	// Check the command to react and answer
	message := strings.ToLower(mes.Content)
	switch message {
	case "!ping":
		b.pingReporter(mes)
	case "!roll":
		b.rollReporter(mes)
	case "!johncena":
		b.jcReporter(mes)
	case "!logs":
		b.logReporter(mes)
	case "!help":
		b.helpReporter(mes)
	case "!boobs":
		b.boobsReporter(mes)
	case "!!terminate":
		panic("Terminating...")
	default:
		b.reactToCommand(mes)
	}
}

func (b *Bot) reactToCommand(mes *discordgo.MessageCreate) {
	// Check the command to react and answer
	message := strings.ToLower(mes.Content)
	command := strings.Split(message, " ")[0]
	switch command {
	case "!status":
		b.statusReporter(mes)
	case "!simcptr":
		b.simcReporter(mes, command, false, true)
	case "!simcstats":
		b.simcReporter(mes, command, true, false)
	case "!simc":
		b.simcReporter(mes, command, false, false)
	case "!queue":
		b.queueReporter(mes)
	case "!realminfo":
		b.realmInfoReporter(mes)
	case "!guildmembers":
		if err := b.SendMessage(mes.ChannelID, m.GuildMembersList); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
		b.guildMembersReporter(mes)
	case "!guildprofs":
		if err := b.SendMessage(mes.ChannelID, m.GuildProfsList); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
		b.guildProfsReporter(mes)
	case "!clean":
		b.cleanUp(mes)
	default:
		if err := b.SendMessage(mes.ChannelID, m.ErrorUser); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
	}
}

/* Tries to call a method and checking if the method returned an error, if it
did check to see if it's HTTP 502 from the Discord API and retry for
`attempts` number of times. */
func retryOnBadGateway(f func() error) {
	var err error
	for i := 0; i < 3; i++ {
		if err = f(); err != nil {
			if strings.HasPrefix(err.Error(), "HTTP 502") {
				// If the error is Bad Gateway, try again after 1 sec.
				time.Sleep(1 * time.Second)
				continue
			} else {
				// Otherwise panic !
				glog.Fatal(err)
			}
		} else {
			// In case of no error, return.
			return
		}
	}
}

func sendFormattedMessage(session *discordgo.Session, chID string, message string) (err error) {
	if i := len(message); len(message) > 1999 {
		for i > 1999 {
			messageSlice := strings.Split(message, "\n")
			mes := messageSlice[0]
			l := len(messageSlice)
			if l == 2 {
				_, err = session.ChannelMessageSend(chID, mes+"\n")
				if err != nil {
					return err
				}
				_, err = session.ChannelMessageSend(chID, messageSlice[1])
				if err != nil {
					return err
				}
				return nil
			}
		Loop:
			for j := 1; j < l-1; j++ {
				mes += "\n" + messageSlice[j]
				if len(mes+"\n"+messageSlice[j+1]) > 1999 {
					if strings.HasPrefix(mes, "```") {
						_, err = session.ChannelMessageSend(chID, mes+"```")
					} else {
						_, err = session.ChannelMessageSend(chID, mes)
					}
					if err != nil {
						return err
					}
					message = strings.Replace(message, mes, "", 1)
					if strings.HasPrefix(mes, "```") {
						message = "```" + message
					}
					break Loop
				}
			}
			i = len(message)
		}
		_, err = session.ChannelMessageSend(chID, message)
	} else {
		_, err = session.ChannelMessageSend(chID, message)
	}
	return
}
