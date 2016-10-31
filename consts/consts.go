package consts

// Admins - users with admin rights
var Admins = []string{
	"208227258954022912", 
	"216928312260296704", 
	"141121531588640769", 
	"217673644854214667",
}

// MatioQuotes - array of quotes of the great Father of all Priests o:)
var MatioQuotes = []string{
	"Если в борделе дела идут плохо, то меняют проституток, а не кровати двигают!",
	"Заебали вы со своими формами!",
	"Лина, у тебя фрейм какой-то не такой, не свой, в общем!",
	"Где мы бля.",
	"Хочется что-то сказать, но тут не добавить, не убавить.",
    "Благословляю тебя, сын мой.",
	"Залупу Вам на эполеты!",
	"Все псы получают рес.",
	"Blizzard х**ни не скажут. Сказали конь - значит конь!",
	"Они ведь это уже говорили! У них синдром Латосия.",
}

const (
	Region          = "eu"
	Locale          = "ru_RU"
	Timezone        = "Europe/Moscow"
	GuildName       = "Аэтернум"
	GuildRealm      = "Ревущий фьорд"
	Pong            = "Pong!"
	Relics          = "https://docs.google.com/spreadsheets/d/11RqT6EIelFWHB1b8f_scFo8sPdXGVYFii_Dr7kkOFLY/edit#gid=1060702296"
	RGB             = "https://docs.google.com/spreadsheets/d/1apphJ2vlZL4eQFZMKeUrYC34PsNt7JFeTZiqNtb0NyE/htmlview?sle=true"
	RealmOn         = "Сервер онлайн! :smile:"
	RealmOff        = "Сервер оффлайн :pensive:"
	RealmHasQueue   = "На сервере очередь, готовься идти делать чай :pensive:"
	RealmHasNoQueue = "Очередей нет, можно заходить! :smile:"
	SpamMessage     = "Напоминаем: поспамив о наборе в гильдию, вы поможете нашему коллективу обрести должное величие среди пантеона богов WoWProgress! :smile:"

	GMCAcquired = "Получаю список согильдейцев из Армори...секундочку :smile:"
	GPCAcquired = "Получаю список профессий в гильдии из Армори...секундочку :smile:"

	Boobies  = "Ах ты грязный извращенец :smile: ну держи)\nhttp://www.gifsfor.com/wp-content/uploads/2012/12/Gifs-for-Tumblr-1445.gif"
	JohnCena = "AND HIS NAME IS JOOOOOOOOOHN CEEEEEEEEEEEENAAAAAAAA! https://youtu.be/QQUgfikLYNI"

	Help = `__**Команды бота:**__

__Общая инфа о гильдии и по прокачке:__
**!roster** - текущий рейдовый состав
**!godbook** - мега-гайд по Легиону
**!relics** - гайдик по реликам на все спеки

**!guildmembers** и **!guildprofs** выводят таблицу состава участников гильдии или профессий в гильдии соответственно.

Список участников можно сортировать по имени (**name**), классу (**class**), спеку(**spec**), уровню (**level**) или уровню предметов (**ilvl**) в порядке возрастания (**asc**) или убывания (**desc**).
Например:
**!guildmembers name=asc level=desc ilvl=desc** выведет список, отсортированный по уровню предметов и уровню в обратном порядке, затем по именам в нормальном порядке.
Этот же порядок применяется по умолчанию, если использовать команду без параметров.

Список профессий можно фильтровать по названию параметром **prof=имя_профессии**
Например:
**!guildprofs prof=Начертание** выведет список всех начертателей в гильдии.

__Команды для WoW'a:__
**!status** ***имя_сервера*** - текущий статус сервера; если не указывать имя - отобразится для РФа
**!queue** ***имя_сервера*** - текущий статус очереди на сервер; если не указывать имя - отобразится для РФа
**!realminfo** ***имя_сервера*** - вся инфа по выбранному серверу; если не указывать имя - отобразится для РФа

С вопросами и предложениями обращаться к **Аэтерису (Илье)**.
Хорошего кача и удачи в борьбе с Легионом! :smile:`

	GuildRosterMID = "218849158721830912"

	WoWAPIRealmsLink         = "https://%v.api.battle.net/wow/realm/status?locale=%v&apikey=%v"
	WoWAPIGuildMembersLink   = "https://%v.api.battle.net/wow/guild/%v/%v?fields=members&locale=%v&apikey=%v"
	WoWAPIGuildNewsLink      = "https://%v.api.battle.net/wow/guild/%v/%v?fields=news&locale=%v&apikey=%v"
	WoWAPICharacterItemsLink = "https://%v.api.battle.net/wow/character/%v/%v?fields=items&locale=%v&apikey=%v"
	WoWAPICharacterProfsLink = "https://%v.api.battle.net/wow/character/%v/%v?fields=professions&locale=%v&apikey=%v"
	WoWAPIItemLink           = "https://%v.api.battle.net/wow/item/%s?locale=%v&apikey=%v"

	WoWArmoryLink     = "http://%v.battle.net/wow/%v/character/%v/%v/advanced"
	WoWArmoryProfLink = "http://%v.battle.net/wow/%v/character/%v/%v/profession/%v"

	WowheadItemLink = "http://ru.wowhead.com/item=%s"

	GoogleAPIShortenerLink = "https://www.googleapis.com/urlshortener/v1/url?key=%v"

	GuildMembersBucketKey = "Guild"
)
