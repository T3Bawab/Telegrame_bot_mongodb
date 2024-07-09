package bot_settings

import (
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	Run  bool
	Bot  *tele.Bot
	Menu *tele.ReplyMarkup
	Ctx  tele.Context
}

var (
	Settings = tele.Settings{
		Token: "7369268252:AAF-1I6jmWWG-FyzeIsBnpjJ6itnNvImZIw",
	}
)
