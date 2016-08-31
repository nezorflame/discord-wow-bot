package main

import (
    "flag"
    "strings"
    "log"
    "time"
    "os"
    "github.com/bwmarrin/discordgo"
    "github.com/nezorflame/discord-wow-bot/wow"
)

const (
    Pong            = "Pong!"
    JohnCena        = "AND HIS NAME IS JOOOOOOOOOHN CEEEEEEEEEEEENAAAAAAAA! https://youtu.be/QQUgfikLYNI"
    Relics          = "https://docs.google.com/spreadsheets/d/11RqT6EIelFWHB1b8f_scFo8sPdXGVYFii_Dr7kkOFLY/edit#gid=1060702296"
    RGB             = "https://docs.google.com/spreadsheets/d/1apphJ2vlZL4eQFZMKeUrYC34PsNt7JFeTZiqNtb0NyE/htmlview?sle=true"
    RealmOn         = "Сервер онлайн! :smile:"
    RealmOff        = "Сервер оффлайн :pensive:"
    RealmHasQueue   = "На сервере очередь, готовься идти делать чай :pensive:"
    RealmHasNoQueue = "Очередей нет, можно заходить! :smile:"

    GuildRosterMID  = "218849158721830912"
)

var (
    logger              *log.Logger
    discordToken        string
    wowToken            string
    mainChannelID       string
    botID               string
)

func init() {
    // Create initials.
	logger = log.New(os.Stderr, "  ", log.Ldate|log.Ltime)

    // Parse command line arguments.
    flag.StringVar(&discordToken, "dt", "", "Account Token")
    flag.StringVar(&wowToken, "wt", "", "WoWAPI dev.battle.net Token")
    flag.StringVar(&mainChannelID, "mc", "", "Main Channel ID")
    flag.Parse()
    if discordToken == "" || wowToken == "" || mainChannelID == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

    wow.InitializeWoWAPI(&wowToken)
}

func logDebug(v ...interface{}) {
	logger.SetPrefix("DEBUG ")
	logger.Println(v...)
}

func logInfo(v ...interface{}) {
	logger.SetPrefix("INFO  ")
	logger.Println(v...)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
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
				panicOnErr(err)
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
		_, err := session.ChannelMessageSend(chID, message)
		return err
	})
    return nil
}

func logPinnedMessages(s *discordgo.Session) {
    logInfo("getPinnedMessages called")
    pinned, err := s.ChannelMessagesPinned(mainChannelID)
    panicOnErr(err)
    logInfo(len(pinned), "messages are pinned:")
    for _, message := range pinned {
        logInfo("[" + message.ID + "]", message.Content)
    }
}

func printMessageByID(s *discordgo.Session, chID string, mesID string) {
    logInfo("printMessageByID called")
    message, err := s.ChannelMessage(mainChannelID, mesID)
    if err != nil {
        logInfo("printMessageByID error: ", err)
        return
    }
    err = sendMessage(s, chID, message.Content)
    panicOnErr(err)
}

func main() {
    logInfo("Logging in...")
    session, err := discordgo.New(discordToken)
    logInfo("Using bot account token...")
    u, err := session.User("@me")
    panicOnErr(err)
    botID = u.ID
    logInfo("Got BotID =", botID)
    setupHandlers(session)
	panicOnErr(err)
    logInfo("Opening session...")
	err = session.Open()
	panicOnErr(err)
	logInfo("Bot is now running.\nPress CTRL-C to exit...")
	<-make(chan struct{})
	return
}

func setupHandlers(session *discordgo.Session) {
	logInfo("Setting up event handlers...")
	session.AddHandler(messageCreate)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == botID {
		return
	}
    // Check the command to answer
    if strings.HasPrefix(m.Content, "!status") {
        statusReporter(s, m)
    }
    if strings.HasPrefix(m.Content, "!queue") {
        queueReporter(s, m)
    }
    if strings.HasPrefix(m.Content, "!realminfo") {
        realmInfoReporter(s, m)
    }
    switch m.Content {
        case "!ping":
            err := sendMessage(s, m.ChannelID, Pong)
            panicOnErr(err)
        case "!johncena":
            err := sendMessage(s, m.ChannelID, JohnCena)
            panicOnErr(err)
        case "!relics":
            err := sendMessage(s, m.ChannelID, Relics)
            panicOnErr(err)
        case "!godbook":
            err := sendMessage(s, m.ChannelID, RGB)
            panicOnErr(err)
        case "!roster":
            printMessageByID(s, m.ChannelID, GuildRosterMID)
        case "!help", "!помощь":
            helpReporter(s, m)
        case "!!printpinned":
            logPinnedMessages(s)
        case "!!terminate":
            panic("Terminating bot...")
        case "!boobs":
            err := sendMessage(s, m.ChannelID, "Покажи фанатам сиськи! :smile:\nhttps://giphy.com/gifs/gene-wilder-z88aYORoi8fQc ")
            panicOnErr(err)
    }
}

func helpReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
    logInfo("Sending help to user...")

    help := "__**Команды бота:**__\n\n"
    help += "__Общая инфа для прокачки и рейдов:__\n"
    help += "**!roster** - текущий рейдовый состав\n"
    help += "**!godbook** - мега-гайд по Легиону\n"
    help += "**!relics** - гайдик по реликам на все спеки\n\n"
    help += "__Команды для WoW'a:__\n"
    help += "**!status** ***имя_сервера*** - текущий статус сервера; если не указывать имя - отобразится для РФа\n"
    help += "**!queue** ***имя_сервера*** - текущий статус очереди на сервер; если не указывать имя - отобразится для РФа\n"
    help += "**!realminfo** ***имя_сервера*** - вся инфа по выбранному серверу; если не указывать имя - отобразится для РФа\n\n"    
    help += "С вопросами и предложениями обращаться к **Аэтерису (Илье)**.\n"
    help += "Хорошего кача и удачи в борьбе с Легионом! :smile:"

    err := sendMessage(s, m.ChannelID, help)
    panicOnErr(err)
}    

func statusReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
    logInfo("getting realm name string...")
    realmString := wow.GetRealmName(m.Content, "!status")
    logInfo(realmString)
    logInfo("getting realm status and sending it...")
    realmStatus, err := wow.GetWoWRealmStatus(realmString)
    if err != nil {
        sendMessage(s, m.ChannelID, err.Error())
    } else if realmStatus {
        err := sendMessage(s, m.ChannelID, RealmOn)
        panicOnErr(err)
    } else {
        err := sendMessage(s, m.ChannelID, RealmOff)
        panicOnErr(err)
    }
}

func queueReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
    logInfo("getting realm name string...")
    realmString := wow.GetRealmName(m.Content, "!queue")
    logInfo(realmString)
    logInfo("getting realm queue status and sending it...")
    realmQueue, err := wow.GetWoWRealmQueueStatus(realmString)
    if err != nil {
        sendMessage(s, m.ChannelID, err.Error())
    } else if realmQueue {
        err := sendMessage(s, m.ChannelID, RealmHasQueue)
        panicOnErr(err)
    } else {
        err := sendMessage(s, m.ChannelID, RealmHasNoQueue)
        panicOnErr(err)
    }
}

func realmInfoReporter(s *discordgo.Session, m *discordgo.MessageCreate) {
    logInfo("getting realm name string...")
    realmString := wow.GetRealmName(m.Content, "!realminfo")
    logInfo(realmString)
    logInfo("getting realm info and sending it...")
    realmInfo, err := wow.GetWoWRealmInfo(realmString)
    if err != nil {
        sendMessage(s, m.ChannelID, err.Error())
    } else {
        err := sendMessage(s, m.ChannelID, realmInfo)
        panicOnErr(err)
    }
}

func containsUser(users []*discordgo.User, userID string) bool {
    for _, u := range users {
        if u.ID == userID {
            return true
        }
    }
    return false
}
