package main

import (
	"log"

	cmd "19u4n4/roebot/commands"
	s "19u4n4/roebot/state"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := t.NewBotAPI("1185985324:AAHuOeP1g9PGgd_cuJno40uGAaKH_nWx0Ew")
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true
	log.Printf("Connected %s", bot.Self.UserName)

	ch, err := bot.GetUpdatesChan(makeUpdateConfig(60))
	if err != nil {
		log.Panic(err)
	}

	transition := cmd.DefaultTransition
	sync := false
	for {
		if sync {
			sync := synchronizer{bot: bot, templates: s.Templates}
			sync.start()
		}
		select {
		case update := <-ch:
			chatID := update.Message.Chat.ID
			hdl := transition(update.Message)
			// username := update.Message.From.UserName
			var reply string
			transition, reply, sync = hdl.Handle()
			if reply != "" {
				msg := t.NewMessage(chatID, reply)
				bot.Send(msg)
			}
		}
	}

}

type synchronizer struct {
	bot       *t.BotAPI
	templates []s.Template
}

func (sync synchronizer) start() {
	sync.pushTemplates()
}

func (sync synchronizer) pushTemplates() {
	for _, tpl := range sync.templates {
		chName := "@" + tpl.TargetChannel
		chat, err := sync.bot.GetChat(t.ChatConfig{SuperGroupUsername: chName})
		if err != nil {
			log.Println(err)
			continue
		}
		msg := t.NewMessage(chat.ID, tpl.Text)
		sync.bot.Send(msg)
	}
}

func makeUpdateConfig(timeout int) t.UpdateConfig {
	cfg := t.NewUpdate(0)
	cfg.Timeout = timeout
	return cfg
}
