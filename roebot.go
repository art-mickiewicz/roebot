package main

import (
	"log"

	cmd "19u4n4/roebot/commands"
	srv "19u4n4/roebot/services"
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
	s.LoadTemplates()
	srv.SyncCBR() // FIXME

	transition := cmd.DefaultTransition
	sync := false
	for {
		if sync {
			sr := synchronizer{bot: bot, templates: s.GetTemplates()}
			sr.start()
		}
		select {
		case update := <-ch:
			if update.EditedMessage != nil {
				chatID := update.EditedMessage.Chat.ID
				msgID := update.EditedMessage.MessageID
				if tpl, ok := s.GetTemplateBySource(s.MessagePtr{ChatID: chatID, MessageID: msgID}); ok {
					tpl.Text = update.EditedMessage.Text
					s.SetTemplate(tpl)
					sync = true
				}
				break
			}
			if update.Message == nil {
				continue
			}
			chatID := update.Message.Chat.ID
			hdl := transition(update.Message)
			// username := update.Message.From.UserName
			var r cmd.Replier
			transition, r, sync = hdl.Handle()
			reply := r.Reply(bot)
			if reply != "" {
				msg := t.NewMessage(chatID, reply)
				msg.ParseMode = "markdown"
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
	s.PersistTemplates()
}

func (sync synchronizer) getChatByName(name string) (t.Chat, error) {
	chName := "@" + name
	return sync.bot.GetChat(t.ChatConfig{SuperGroupUsername: chName})
}

func (sync synchronizer) pushTemplates() {
	for i, tpl := range sync.templates {
		if tpl.IsPosted() {
			if tpl.TargetMessagePtr.ChatID == 0 {
				chat, err := sync.getChatByName(tpl.TargetChannel)
				if err != nil {
					log.Println(err)
					continue
				}
				tpl.TargetMessagePtr.ChatID = chat.ID
				sync.templates[i] = tpl
			}
			edit := t.EditMessageTextConfig{
				BaseEdit: t.BaseEdit{
					ChatID:    tpl.TargetMessagePtr.ChatID,
					MessageID: tpl.TargetMessagePtr.MessageID,
				},
				Text: tpl.Text,
			}
			sync.bot.Send(edit)
		} else {
			chat, err := sync.getChatByName(tpl.TargetChannel)
			if err != nil {
				log.Println(err)
				continue
			}
			msg := t.NewMessage(chat.ID, tpl.Text)
			postedMsg, err := sync.bot.Send(msg)
			if err == nil {
				tpl.TargetMessagePtr = s.MessagePtr{ChatID: chat.ID, MessageID: postedMsg.MessageID}
				sync.templates[i] = tpl
			}
		}
	}
}

func makeUpdateConfig(timeout int) t.UpdateConfig {
	cfg := t.NewUpdate(0)
	cfg.Timeout = timeout
	return cfg
}
