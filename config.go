package main

import (
	rate "github.com/juju/ratelimit"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	o = &Options{}
	m = &Messages{}
	// RateLimiter for GET requests to WoW API
	RateLimiter *rate.Bucket
)

// LoadConfig loads the config from the config.toml
func LoadConfig(sl *zap.SugaredLogger) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/opt/bot")

	sl.Info("Loading configuration...")
	if err := viper.ReadInConfig(); err != nil {
		sl.Fatalf("Unable to load config: %s", err)
	}

	// Mandatory tokens
	if o.DiscordToken = viper.GetString("discord_token"); o.DiscordToken == "" {
		sl.Fatal("'discord_token' must be present")
	}
	if o.WoWToken = viper.GetString("wow_token"); o.WoWToken == "" {
		sl.Fatal("'wow_token' must be present")
	}
	if o.GoogleToken = viper.GetString("google_token"); o.GoogleToken == "" {
		sl.Fatal("'google_token' must be present")
	}

	// Links
	o.WowheadItemLink = viper.GetString("wowhead_item")
	o.WoWDBItemLink = viper.GetString("wowdb_item")
	o.GoogleShortenerLink = viper.GetString("google_shortener")

	// Legendary check period
	o.CharacterCheckPeriod = viper.GetDuration("character_check_period")
	o.LegendaryCheckPeriod = viper.GetDuration("legendary_check_period")
	o.LegendaryRelevancePeriod = viper.GetDuration("legendary_rel_period")

	// SimC
	o.SimcDir = viper.GetString("simc.dir")
	o.SimcCmdStable = viper.GetString("simc.cmd_stable")
	o.SimcCmdPtr = viper.GetString("simc.cmd_ptr")
	o.SimcArgsImport = viper.GetString("simc.args_import")
	o.SimcArgsNoStats = viper.GetString("simc.args_no_stats")
	o.SimcArgsWithStats = viper.GetString("simc.args_with_stats")

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
	if o.APIRateLimit = viper.GetInt64("blizzard.api_rate_limit"); o.APIRateLimit <= 0 {
		sl.Fatal("'blizzard.api_rate_limit' must be > 0")
	}
	if o.APIMaxRetries = viper.GetInt("blizzard.api_max_retries"); o.APIMaxRetries <= 0 {
		sl.Fatal("'blizzard.api_max_retries' must be > 0")
	}
	RateLimiter = rate.NewBucketWithRate(float64(o.APIRateLimit), o.APIRateLimit)

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
	m.BoobiesPrefix = viper.GetString("messages.boobies_prefix")
	m.Boobies = viper.GetStringSlice("messages.boobies")
	m.JohnCena = viper.GetString("messages.johncena")
	m.Legendary = viper.GetString("messages.legendary")
	m.Logs = viper.GetString("messages.logs")

	m.Roll1 = viper.GetString("messages.roll_1")
	m.RollX = viper.GetString("messages.roll_x")
	m.Roll100 = viper.GetString("messages.roll_100")

	m.RealmInfo = viper.GetString("messages.realm_info")
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

	sl.Info("Configuration is loaded successfully")
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
