package main

import (
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

var (
	o = &Options{}
	m = &Messages{}
)

// LoadConfig loads the config from the config.toml
func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	glog.Info("Loading configuration...")
	if err := viper.ReadInConfig(); err != nil {
		glog.Fatalf("Unable to load config: %s", err)
	}

	// Mandatory tokens
	if o.DiscordToken = viper.GetString("discord_token"); o.DiscordToken == "" {
		glog.Fatal("'discord_token' must be present")
	}
	if o.WoWToken = viper.GetString("wow_token"); o.WoWToken == "" {
		glog.Fatal("'wow_token' must be present")
	}
	if o.GoogleToken = viper.GetString("google_token"); o.GoogleToken == "" {
		glog.Fatal("'google_token' must be present")
	}
	if o.Bucket = viper.GetString("bucket"); o.Bucket == "" {
		glog.Fatal("'bucket' must be present")
	}

	// SSH
	if o.SSHAddress = viper.GetString("ssh.addr"); o.SSHAddress == "" {
		glog.Fatal("'ssh.addr' must be present")
	}
	if o.SSHUser = viper.GetString("ssh.user"); o.SSHUser == "" {
		glog.Fatal("'ssh.user' must be present")
	}

	// Links
	o.WowheadItemLink = viper.GetString("wowhead_item")
	o.WoWDBItemLink = viper.GetString("wowdb_item")
	o.GoogleShortenerLink = viper.GetString("google_shortener")

	// SimC
	o.SimcImport = viper.GetString("simc.import")
	o.SimcNoStats = viper.GetString("simc.no_stats")
	o.SimcWithStats = viper.GetString("simc.with_stats")

	// Discord
	o.Admins = viper.GetStringSlice("discord.admins")
	o.GeneralChannelID = viper.GetString("discord.general_channel")

	// Guild
	o.GuildRegion = viper.GetString("guild.region")
	o.GuildLocale = viper.GetString("guild.locale")
	o.GuildTimezone = viper.GetString("guild.timezone")
	o.GuildName = viper.GetString("guild.name")
	o.GuildRealm = viper.GetString("guild.realm")

	// WoW slices
	o.WoWClasses = getMapWithNames("classes")
	o.WoWGenders = getMapWithNames("genders")
	o.WoWFactions = getMapWithNames("factions")
	o.WoWRaces = getMapWithNames("races")
	o.WoWProfessions = getMapWithNames("professions")

	// Blizzard
	o.APIRealmsLink = viper.GetString("blizzard.api_realms")
	o.APIGuildMembersLink = viper.GetString("blizzard.api_guild_members")
	o.APIGuildNewsLink = viper.GetString("blizzard.api_guild_news")
	o.APICharItemsLink = viper.GetString("blizzard.api_char_items")
	o.APICharNewsLink = viper.GetString("blizzard.api_char_news")
	o.APICharProfsLink = viper.GetString("blizzard.api_char_profs")
	o.APIItemLink = viper.GetString("blizzard.api_item")

	o.ArmoryCharLink = viper.GetString("blizzard.armory_char")
	o.ArmoryProfLink = viper.GetString("blizzard.armory_prof")

	// Messages
	m.ErrorUser = viper.GetString("messages.error_user")
	m.ErrorServer = viper.GetString("messages.error_server")
	m.Help = viper.GetString("messages.help")
	m.Pong = viper.GetString("messages.pong")
	m.Clean = viper.GetString("messages.clean")
	m.Boobies = viper.GetString("messages.boobies")
	m.JohnCena = viper.GetString("messages.johncena")
	m.Legendary = viper.GetString("messages.legendary")
	m.Logs = viper.GetString("messages.logs")

	m.Roll1 = viper.GetString("messages.roll_1")
	m.RollX = viper.GetString("messages.roll_x")
	m.Roll100 = viper.GetString("messages.roll_100")

	m.Relics = viper.GetString("messages.relics")
	m.Godbook = viper.GetString("messages.godbook")

	m.RealmOn = viper.GetString("messages.realm_on")
	m.RealmOff = viper.GetString("messages.realm_off")
	m.RealmQueue = viper.GetString("messages.realm_queue")
	m.RealmNoQueue = viper.GetString("messages.realm_noqueue")

	m.GuildMembersList = viper.GetString("messages.guild_members_list")
	m.GuildProfsList = viper.GetString("messages.guild_profs_list")

	m.SimcArmory = viper.GetString("messages.simc_armory")
	m.SimcArmoryError = viper.GetString("messages.simc_armory_error")
	m.SimcImport = viper.GetString("messages.simc_import")
	m.SimcImportSuccess = viper.GetString("messages.simc_import_success")
	m.SimcProfile = viper.GetString("messages.simc_profile")

	glog.Info("Configuration is loaded successfully")
}

func getMapWithNames(confName string) (confMap map[int64]string) {
	confMap = make(map[int64]string)
	confSlice := viper.Get(confName).([]interface{})
	for _, i := range confSlice {
		id := i.(map[string]interface{})["id"].(int64)
		name := i.(map[string]interface{})["name"].(string)
		confMap[id] = name
	}
	return
}
