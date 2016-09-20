package consts

const (
    Region                      = "eu"
    Locale                      = "ru_RU"
    GuildName                   = "Аэтернум"
    GuildRealm                  = "Ревущий фьорд"
    Pong                        = "Pong!"
    Relics                      = "https://docs.google.com/spreadsheets/d/11RqT6EIelFWHB1b8f_scFo8sPdXGVYFii_Dr7kkOFLY/edit#gid=1060702296"
    RGB                         = "https://docs.google.com/spreadsheets/d/1apphJ2vlZL4eQFZMKeUrYC34PsNt7JFeTZiqNtb0NyE/htmlview?sle=true"
    RealmOn                     = "Сервер онлайн! :smile:"
    RealmOff                    = "Сервер оффлайн :pensive:"
    RealmHasQueue               = "На сервере очередь, готовься идти делать чай :pensive:"
    RealmHasNoQueue             = "Очередей нет, можно заходить! :smile:"

    GMCAcquired                 = "Получаю список согильдейцев из Армори...подожди секунд 10-15 :smile:"
    GPCAcquired                 = "Получаю список профессий в гильдии из Армори...подожди секунд 10-15 :smile:"

    Boobies                     = "Покажи фанатам сиськи! :smile:\nhttps://giphy.com/gifs/gene-wilder-z88aYORoi8fQc"
    JohnCena                    = "AND HIS NAME IS JOOOOOOOOOHN CEEEEEEEEEEEENAAAAAAAA! https://youtu.be/QQUgfikLYNI"

    Help                        = `__**Команды бота:**__

__Общая инфа о гильдии и по прокачке:__
**!guildmembers** - состав гильдии
**!guildprofs** - список всех профессий в гильдии
**!roster** - текущий рейдовый состав
**!godbook** - мега-гайд по Легиону
**!relics** - гайдик по реликам на все спеки

__Команды для WoW'a:__
**!status** ***имя_сервера*** - текущий статус сервера; если не указывать имя - отобразится для РФа
**!queue** ***имя_сервера*** - текущий статус очереди на сервер; если не указывать имя - отобразится для РФа
**!realminfo** ***имя_сервера*** - вся инфа по выбранному серверу; если не указывать имя - отобразится для РФа

С вопросами и предложениями обращаться к **Аэтерису (Илье)**.
Хорошего кача и удачи в борьбе с Легионом! :smile:`

    GuildRosterMID              = "218849158721830912"

    WoWAPIRealmsLink            = "https://%v.api.battle.net/wow/realm/status?locale=%v&apikey=%v"
    WoWAPIGuildMembersLink      = "https://%v.api.battle.net/wow/guild/%v/%v?fields=members&locale=%v&apikey=%v"
    WoWAPIGuildNewsLink         = "https://%v.api.battle.net/wow/guild/%v/%v?fields=news&locale=%v&apikey=%v"
    WoWAPICharacterItemsLink    = "https://%v.api.battle.net/wow/character/%v/%v?fields=items&locale=%v&apikey=%v"
    WoWAPICharacterProfsLink    = "https://%v.api.battle.net/wow/character/%v/%v?fields=professions&locale=%v&apikey=%v"
    WoWAPIItemLink              = "https://%v.api.battle.net/wow/item/%d?locale=%v&apikey=%v"

    WoWArmoryProfLink           = "http://%v.battle.net/wow/%v/character/%v/%v/profession/%v"

    WowheadItemLink             = "http://ru.wowhead.com/item=%d"

    GoogleAPIShortenerLink      = "https://www.googleapis.com/urlshortener/v1/url?key=%v"
)