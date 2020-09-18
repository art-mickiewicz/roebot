package commands

import t "github.com/go-telegram-bot-api/telegram-bot-api"

type Replier interface {
	Reply(bot *t.BotAPI) string
}

type str string

func (s str) Reply(bot *t.BotAPI) string {
	return string(s)
}

type funcReplier func(*t.BotAPI) string

func (r funcReplier) Reply(bot *t.BotAPI) string {
	return r(bot)
}

type Transition func(*t.Message) Handler

type Handler interface {
	Handle() (transitTo Transition, r Replier, sync bool)
}
