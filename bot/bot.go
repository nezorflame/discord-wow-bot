package bot

import (
	"log"
	"strings"
    "strconv"
	"time"

	"github.com/arteev/fmttab"
	"github.com/bwmarrin/discordgo"
	"github.com/nezorflame/discord-wow-bot/consts"
	"github.com/nezorflame/discord-wow-bot/wow"
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

    // Users - map of guild members
    Users map[string]string

	botID   string
	session *discordgo.Session
)

// Start - function to start Discord bot
func Start() {
	// Fix for a new Discord Bot token auth
	DiscordToken = "Bot " + DiscordToken
	wow.InitializeWoWAPI(&WoWToken, &GoogleToken)
	logInfo("Logging in...")
	session, err := discordgo.New(DiscordToken)
	logInfo("Using bot account token...")
	u, err := session.User("@me")
	logOnErr(err)
	botID = u.ID
	logInfo("Got BotID =", botID)
	logInfo("Adding handlers...")
	setup(session)
	logInfo("Opening session...")
	err = session.Open()
	logOnErr(err)
	logInfo("Starting guild watcher and spammer...")
	go RunGuildWatcher(session)
	go RunGuildSpammer(session)
}

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
		Logger.Println(err)
	}
}

/* Tries to call a method and checking if the method returned an error, if it
did check to see if it's HTTP 502 from the Discord API and retry for
`attempts` number of times. */
func retryOnBadGateway(f func() error) {
	var err error
	for i := 0; i < 3; i++ {
		err = f()
		if err != nil {
			if strings.HasPrefix(err.Error(), "HTTP 502") {
				// If the error is Bad Gateway, try again after 1 sec.
				time.Sleep(1 * time.Second)
				continue
			} else {
				// Otherwise panic !
				logOnErr(err)
			}
		} else {
			// In case of no error, return.
			return
		}
	}
}

func sendMessage(session *discordgo.Session, chID string, message string) error {
	logInfo("SENDING MESSAGE:", message)
	retryOnBadGateway(func() error {
		err := sendFormattedMessage(session, chID, message)
		return err
	})
	return nil
}

func sendFormattedMessage(session *discordgo.Session, chID string, fullMessage string) error {
	var err error
	message := fullMessage
	i := len(message)
	if len(message) > 1999 {
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
		if err != nil {
			return err
		}
	} else {
		_, err = session.ChannelMessageSend(chID, fullMessage)
	}
	return err
}

func logPinnedMessages(s *discordgo.Session) {
	logInfo("getPinnedMessages called")
	pinned, err := s.ChannelMessagesPinned(DiscordMChanID)
	logOnErr(err)
	logInfo(len(pinned), "messages are pinned:")
	for _, message := range pinned {
		logInfo("["+message.ID+"]", message.Content)
	}
}

func printMessageByID(s *discordgo.Session, chID string, mesID string) {
	logInfo("printMessageByID called")
	message, err := s.ChannelMessage(DiscordMChanID, mesID)
	if err != nil {
		logInfo("printMessageByID error: ", err)
		return
	}
	err = sendMessage(s, chID, message.Content)
	logOnErr(err)
}

func setup(session *discordgo.Session) {
	logInfo("Setting up event handlers...")
	session.AddHandler(messageCreate)
}

// RunGuildWatcher - function for starting the guild news watcher
func RunGuildWatcher(s *discordgo.Session) {
	// TODO: Very dirty, need to rewrite
	legendaries := make(map[string]bool)
	for {
		messages, err := wow.GetGuildLegendaries(consts.GuildRealm, consts.GuildName)
		panicOnErr(err)
		for _, m := range messages {
			if _, ok := legendaries[m]; !ok {
				err := sendMessage(s, DiscordMChanID, m)
				logOnErr(err)
				log.Println(m)
				legendaries[m] = true
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

// RunGuildSpammer - function for SPAMMING :)
func RunGuildSpammer(s *discordgo.Session) {
	for {
		time.Sleep(2 * time.Hour)
		if timeIsAllowed() {
            err := sendMessage(s, DiscordMChanID, consts.SpamMessage)
			logOnErr(err)
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == botID {
		return
	}
	// Check the command to react and answer
	if strings.HasPrefix(m.Content, "!status") {
		statusReporter(s, m)
	}
	if strings.HasPrefix(m.Content, "!queue") {
		queueReporter(s, m)
	}
	if strings.HasPrefix(m.Content, "!realminfo") {
		realmInfoReporter(s, m)
	}
	if strings.HasPrefix(m.Content, "!guildmembers") {
		err := sendMessage(s, m.ChannelID, consts.GMCAcquired)
		logOnErr(err)
		guildMembersReporter(s, m)
	}
	if strings.HasPrefix(m.Content, "!guildprofs") {
		err := sendMessage(s, m.ChannelID, consts.GPCAcquired)
		logOnErr(err)
		guildProfsReporter(s, m)
	}
	if strings.HasPrefix(m.Content, "!clean") {
		cleanUp(s, m)
	}
	switch m.Content {
	case "!ping":
		err := sendMessage(s, m.ChannelID, consts.Pong)
		logOnErr(err)
	case "!johncena":
		err := sendMessage(s, m.ChannelID, consts.JohnCena)
		logOnErr(err)
	case "!relics":
		err := sendMessage(s, m.ChannelID, consts.Relics)
		logOnErr(err)
	case "!godbook":
		err := sendMessage(s, m.ChannelID, consts.RGB)
		logOnErr(err)
	case "!roster":
		printMessageByID(s, m.ChannelID, consts.GuildRosterMID)
	case "!help", "!помощь":
		helpReporter(s, m)
	case "!boobs":
		boobsReporter(s, m)
	case "!!printpinned":
		logPinnedMessages(s)
	case "!!terminate":
		panic("Terminating bot...")
	}
}

func cleanUp(s *discordgo.Session, m *discordgo.MessageCreate) {
	logInfo("Removing bot messages...")
    var err error
	user := m.Author.Username
    if m.ChannelID == DiscordMChanID && !containsUser(consts.Admins, m.Author.ID) {
        logInfo("User is trying to delete all bot messages from main channel! Won't work!")
        err := sendMessage(s, m.ChannelID, "Прости, но в главном чате мои сообщения могут удалять только админы :smile:")
		logOnErr(err)
        return
    }
    am := strings.Replace(m.Message.Content, "!clean", "", 1)
    am = strings.Replace(am, " ", "", -1)
    logInfo("User", user, "- amount to delete:", am)
    var amount int
    switch am {
        case "all":
            amount = -1
        case "":
            amount = 1
        default:
            amount, err = strconv.Atoi(am)
            if err != nil {
                logOnErr(err)
                return
            }
    }
    lastMessageChecked := m.ID
    chanMessages, _ := s.ChannelMessages(m.ChannelID, 100, lastMessageChecked, "")
    mesToDelete := make(map[string]string)
    for {
        if len(mesToDelete) == amount {
            break
        }
        for _, mes := range chanMessages {
            logInfo(mes.ID, mes.Author.Username, mes.Author.ID)
            lastMessageChecked = mes.ID
            if mes.Author.ID == botID {
                if _, ok := mesToDelete[mes.ID]; !ok {
                    mesToDelete[mes.ID] = mes.ID
                }
                if len(mesToDelete) == amount {
                    break
                }
            }
        }
        chm, _ := s.ChannelMessages(m.ChannelID, 100, lastMessageChecked, "")
        if compareMesArrays(chm, chanMessages) {
            logInfo("Reached the end, exiting loop...")
            break
        }
        chanMessages = chm
    }
    for _, mID := range mesToDelete {
        err = s.ChannelMessageDelete(m.ChannelID, mID)
        logOnErr(err)
    }
    logInfo("Deleted all messages")
    return
}

func helpReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
	logInfo("Sending help to user...")
	err := sendMessage(s, m.ChannelID, consts.Help)
	logOnErr(err)
}

func boobsReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
	logInfo("Sending boobies to user...:)")
	err := sendMessage(s, m.ChannelID, consts.Boobies)
	logOnErr(err)
}

func statusReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
	logInfo("getting realm name string...")
	realmString := wow.GetRealmName(m.Content, "!status")
	logInfo(realmString)
	logInfo("getting realm status and sending it...")
	realmStatus, err := wow.GetRealmStatus(realmString)
	if err != nil {
		sendMessage(s, m.ChannelID, err.Error())
	} else if realmStatus {
		err := sendMessage(s, m.ChannelID, consts.RealmOn)
		logOnErr(err)
	} else {
		err := sendMessage(s, m.ChannelID, consts.RealmOff)
		logOnErr(err)
	}
}

func queueReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
	logInfo("getting realm name string...")
	realmString := wow.GetRealmName(m.Content, "!queue")
	logInfo(realmString)
	logInfo("getting realm queue status and sending it...")
	realmQueue, err := wow.GetRealmQueueStatus(realmString)
	if err != nil {
		sendMessage(s, m.ChannelID, err.Error())
	} else if realmQueue {
		err := sendMessage(s, m.ChannelID, consts.RealmHasQueue)
		logOnErr(err)
	} else {
		err := sendMessage(s, m.ChannelID, consts.RealmHasNoQueue)
		logOnErr(err)
	}
}

func guildMembersReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
	logInfo("getting parametes string slice...")
	var parameters []string
	realmString := consts.GuildRealm
	guildNameString := consts.GuildName
	paramString := strings.TrimPrefix(m.Content, "!guildmembers")
	paramString = strings.TrimPrefix(paramString, " ")
	if paramString != "" {
		parameters = strings.Split(paramString, " ")
		logInfo("paramString:", paramString, "parameters len:", len(parameters))
	}
	logInfo("getting guild members list and sending it...")
	guildMembersInfo, err := wow.GetGuildMembers(realmString, guildNameString, parameters)
	if err != nil {
		sendMessage(s, m.ChannelID, err.Error())
		return
	}
	tab := fmttab.New("Список согильдейцев", fmttab.BorderDouble, nil)
	tab.AddColumn("Имя", fmttab.WidthAuto, fmttab.AlignLeft).
		AddColumn("Уровень", 7, fmttab.AlignLeft).
		AddColumn("Класс", 18, fmttab.AlignLeft).
		AddColumn("Специализация", 18, fmttab.AlignLeft).
		AddColumn("iLevel", 6, fmttab.AlignLeft).
		AddColumn("Армори", 22, fmttab.AlignLeft)
	for _, member := range guildMembersInfo {
		tab.AppendData(map[string]interface{}{
			"Имя":           member["Name"],
			"Уровень":       member["Level"],
			"Класс":         member["Class"],
			"Специализация": member["Spec"],
			"iLevel":        member["ItemLevel"],
			"Армори":        member["Link"],
		})
	}
	err = sendMessage(s, m.ChannelID, "```"+tab.String()+"```")
	logInfo(len(tab.String()))
	logOnErr(err)
}

func guildProfsReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
	logInfo("getting parametes string slice...")
	realmString := consts.GuildRealm
	guildNameString := consts.GuildName
	paramString := strings.TrimPrefix(m.Content, "!guildprofs")
	paramString = strings.TrimPrefix(paramString, " ")
	logInfo("getting guild profs list and sending it...")
	guildProfsInfo, err := wow.GetGuildProfs(realmString, guildNameString, paramString)
	if err != nil {
		sendMessage(s, m.ChannelID, err.Error())
		return
	}
	tab := fmttab.New("Список профессий в гильдии", fmttab.BorderDouble, nil)
	tab.AddColumn("Имя", fmttab.WidthAuto, fmttab.AlignLeft).
		AddColumn("1 профа", 15, fmttab.AlignLeft).
		AddColumn("Уровень 1 профы", fmttab.WidthAuto, fmttab.AlignLeft).
		AddColumn("2 профа", 15, fmttab.AlignLeft).
		AddColumn("Уровень 2 профы", fmttab.WidthAuto, fmttab.AlignLeft)
	for _, member := range guildProfsInfo {
		tab.AppendData(map[string]interface{}{
			"Имя":             member["Name"],
			"1 профа":         member["FirstProf"],
			"Уровень 1 профы": member["FirstProfLevel"],
			"2 профа":         member["SecondProf"],
			"Уровень 2 профы": member["SecondProfLevel"],
		})
	}
	err = sendMessage(s, m.ChannelID, "```"+tab.String()+"```")
	logOnErr(err)
}

func realmInfoReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
	logInfo("getting realm name string...")
	realmString := wow.GetRealmName(m.Content, "!realminfo")
	logInfo(realmString)
	logInfo("getting realm info and sending it...")
	realmInfo, err := wow.GetRealmInfo(realmString)
	if err != nil {
		sendMessage(s, m.ChannelID, err.Error())
	} else {
		err := sendMessage(s, m.ChannelID, realmInfo)
		logOnErr(err)
	}
}

func containsUser(users []string, userID string) bool {
	for _, u := range users {
		if u == userID {
			return true
		}
	}
	return false
}

func compareMesArrays(a, b []*discordgo.Message) bool {
    for i := range a {
        if a[i].ID != b[i].ID {
            return false
        }
    }
    return true
}

func timeIsAllowed() bool {
    location, err := time.LoadLocation(consts.Timezone)
    panicOnErr(err)
    now := time.Now().In(location)
    hour := now.Hour()
    weekday := now.Weekday()
    switch weekday {
        // saturday has raids and is a holiday
        case time.Saturday:
            if !(hour >= 2 && hour <= 10 || hour >= 20 && hour <= 23) {
                logInfo("Saturday spam :) time now:", now.String())
                return true
            }
        // sunday 
        case time.Sunday:
            if !(hour >= 2 && hour <= 8) {
                logInfo("Sunday spam :) time now:", now.String())
                return true
            }
        // work days
        default:
            if !(hour >= 2 && hour <= 8 || hour >= 20 && hour <= 23) {
                logInfo("Workday spam :) time now:", now.String())
                return true
            }
    }
    return false
}