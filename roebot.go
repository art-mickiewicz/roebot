package main

import (
	"log"

	cmd "19u4n4/roebot/commands"
	"19u4n4/roebot/config"
	srv "19u4n4/roebot/services"
	_ "19u4n4/roebot/services/binance"
	_ "19u4n4/roebot/services/cbr"
	s "19u4n4/roebot/state"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *t.BotAPI // FIXME

func main() {
	var err error
	bot, err = t.NewBotAPI("1185985324:AAHuOeP1g9PGgd_cuJno40uGAaKH_nWx0Ew")
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
	srv.SyncAll()

	transition := cmd.DefaultTransition
	sync := false
	go func() {
		for range srv.Updates {
			sr := synchronizer{bot: bot, templates: s.GetTemplates()}
			sr.start()
		}
	}()

	for {
		if sync {
			sr := synchronizer{bot: bot, templates: s.GetTemplates()}
			sr.start()
		}
		select {
		case update := <-ch:
			if update.EditedMessage != nil {
				if allowed := checkAccess(update.EditedMessage); !allowed {
					accessDeniedMessage(update.EditedMessage.Chat.ID)
					break
				}
				chatID := update.EditedMessage.Chat.ID
				msgID := update.EditedMessage.MessageID
				if tpl, ok := s.GetTemplateBySource(s.MessagePtr{ChatID: chatID, MessageID: msgID}); ok {
					tpl.Text = update.EditedMessage.Text
					if ok := s.SetTemplate(tpl); ok {
						sync = true
					} else {
						msg := t.NewMessage(chatID, "Ошибка в шаблоне.")
						msg.ParseMode = "markdown"
						bot.Send(msg)
					}
				}
				break
			}
			if update.Message == nil {
				break
			}
			if allowed := checkAccess(update.Message); !allowed {
				accessDeniedMessage(update.Message.Chat.ID)
				break
			}
			chatID := update.Message.Chat.ID
			hdl := transition(update.Message)
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

func checkAccess(message *t.Message) bool {
	username := message.From.UserName
	for _, admin := range config.Admins {
		if username == admin {
			return true
		}
	}
	return false
}

func accessDeniedMessage(chatID int64) {
	msg := t.NewMessage(chatID, "Я с тобой не разговариваю, обратись к админам.")
	msg.ParseMode = "markdown"
	bot.Send(msg)
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
				Text: tpl.Apply(srv.GetVariablesValues()),
			}
			edit.ParseMode = "markdown"
			sync.bot.Send(edit)
		} else {
			chat, err := sync.getChatByName(tpl.TargetChannel)
			if err != nil {
				log.Println(err)
				continue
			}
			msg := t.NewMessage(chat.ID, tpl.Apply(srv.GetVariablesValues()))
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
