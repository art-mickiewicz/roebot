package main

import (
	"log"
	"robpike.io/filter"

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
			if update.EditedMessage != nil {
				chatID := update.Message.Chat.ID
				tpl := filter.Choose(s.Templates, func(x Template) bool {
					return x != ""
				}).([]Template)
				hdl := cmd.SetTemplate{TemplateID: }
				break
			}
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
	for i, tpl := range sync.templates {
		chName := "@" + tpl.TargetChannel
		chat, err := sync.bot.GetChat(t.ChatConfig{SuperGroupUsername: chName})
		if err != nil {
			log.Println(err)
			continue
		}
		if tpl.IsPosted() {
			edit := t.EditMessageTextConfig{
				BaseEdit: t.BaseEdit{
					ChatID:    chat.ID,
					MessageID: tpl.TargetMessageID,
				},
				Text: tpl.Text,
			}
			sync.bot.Send(edit)
		} else {
			msg := t.NewMessage(chat.ID, tpl.Text)
			postedMsg, err := sync.bot.Send(msg)
			if err == nil {
				tpl.TargetMessageID = postedMsg.MessageID
				sync.templates[i] = tpl
				log.Println(sync.templates)
			}
		}
	}
}

func makeUpdateConfig(timeout int) t.UpdateConfig {
	cfg := t.NewUpdate(0)
	cfg.Timeout = timeout
	return cfg
}
