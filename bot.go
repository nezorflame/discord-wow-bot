package main

import (
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
)

// Bot struct for the Discord bot
type Bot struct {
	ID      string
	Session *discordgo.Session

	Guild *GuildInfo

	HighLvlCharacters map[string]*Character
	LegendariesByChar map[string][]*Item

	CharMutex sync.Mutex

	SL *zap.SugaredLogger
}

// Start launches the connection to the bot
func (b *Bot) Start() {
	var (
		err error
		u   *discordgo.User
	)

	b.SL.Info("Logging in...")
	if b.Session, err = discordgo.New(o.DiscordToken); err != nil {
		b.SL.Fatalf("Unable to connect to Discord: %s", err)
	}

	b.SL.Info("Using bot account token...")
	if u, err = b.Session.User("@me"); err != nil {
		b.SL.Fatalf("Unable to get @me: %s", err)
	} else {
		b.ID = u.ID
		b.SL.Infof("Got BotID = %s", b.ID)
	}

	b.SL.Info("Adding handlers...")
	b.Session.AddHandler(b.parseMessage)

	b.SL.Info("Opening session...")
	if err = b.Session.Open(); err != nil {
		b.SL.Fatalf("Unable to open the session: %s", err)
	}

	b.LegendariesByChar = make(map[string][]*Item)

	WoWItemsMap = make(map[string]*Item)

	b.SL.Info("Starting guild watcher...")
	go b.guildWatcher()

	// wait a bit for a guild watcher to start
	time.Sleep(time.Second)

	b.SL.Info("Starting legendaries watcher...")
	go b.legendaryWatcher()

	b.SL.Info("Bot started.")
}

// SendMessage sends the message to the selected channel
func (b *Bot) SendMessage(chID string, message string) {
	var err error
	b.SL.Infof("SENDING MESSAGE: %s", message)
	err = retryOnBadGateway(func() error {
		return sendFormattedMessage(b.Session, chID, message)
	})
	if err != nil {
		b.SL.Errorf("Unable to send the message: %s", err)
	}
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (b *Bot) parseMessage(s *discordgo.Session, mes *discordgo.MessageCreate) {
	// Try to work with the file if bot is a recepient
	if len(mes.Attachments) == 1 && userListContainsUser(mes.Mentions, b.ID) {
		b.reactToFile(mes)
	}
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

func (b *Bot) reactToFile(mes *discordgo.MessageCreate) {
	// Ask user how to react
	b.SL.Info("trying to process the file...")
	msgParts := strings.Split(mes.Content, " ")
	if len(msgParts) != 3 {
		b.SL.Infof("Command is incorrect: %s", mes.Content)
		b.SendMessage(mes.ChannelID, m.ErrorUser)
		return
	}
	command := msgParts[1]
	switch command {
	case "!simcptr":
		b.simcProfileReporter(mes, command, false, true)
	case "!simcstats":
		b.simcProfileReporter(mes, command, true, false)
	case "!simc":
		b.simcProfileReporter(mes, command, false, false)
	default:
		b.SL.Infof("Command is incorrect: %s", mes.Content)
		b.SendMessage(mes.ChannelID, m.ErrorUser)
	}

}

func (b *Bot) reactToCommand(mes *discordgo.MessageCreate) {
	// Check the command to react and answer
	command := strings.Split(strings.ToLower(mes.Content), " ")[0]
	switch command {
	case "!status":
		b.statusReporter(mes)
	case "!simcptr":
		b.simcArmoryReporter(mes, command, false, true)
	case "!simcstats":
		b.simcArmoryReporter(mes, command, true, false)
	case "!simc":
		b.simcArmoryReporter(mes, command, false, false)
	case "!queue":
		b.queueReporter(mes)
	case "!realminfo":
		b.realmInfoReporter(mes)
	case "!guildmembers":
		b.guildMembersReporter(mes)
	case "!guildprofs":
		b.guildProfsReporter(mes)
	case "!clean":
		b.cleanUp(mes)
	default:
		b.SL.Infof("Command is incorrect: %s", mes.Content)
		b.SendMessage(mes.ChannelID, m.ErrorUser)
	}
}

/* Tries to call a method and checking if the method returned an error, if it
did check to see if it's HTTP 502 from the Discord API and retry for
`attempts` number of times. */
func retryOnBadGateway(f func() error) error {
	var err error
	for i := 0; i < 3; i++ {
		if err = f(); err != nil {
			if strings.HasPrefix(err.Error(), "HTTP 502") {
				// If the error is Bad Gateway, try again after 1 sec.
				time.Sleep(1 * time.Second)
				continue
			} else {
				// Otherwise return error
				return err
			}
		} else {
			// In case of no error, return nil
			return nil
		}
	}
	return err
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

func userListContainsUser(users []*discordgo.User, userID string) bool {
	for _, u := range users {
		if u.ID == userID {
			return true
		}
	}
	return false
}
