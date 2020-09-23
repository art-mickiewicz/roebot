package main

import (
	"errors"
	"log"
	"strconv"

	cmd "19u4n4/roebot/commands"
	cfg "19u4n4/roebot/config"
	rich "19u4n4/roebot/richtext"
	srv "19u4n4/roebot/services"
	_ "19u4n4/roebot/services/binance"
	_ "19u4n4/roebot/services/cbr"
	s "19u4n4/roebot/state"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *t.BotAPI

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
	s.LoadChats()
	s.LoadTemplates()
	go srv.SyncAll()

	transition := cmd.DefaultTransition
	sync := false

	for {
		if sync {
			Sync()
		}
		select {
		case <-srv.Updates:
			Sync()
		case update := <-ch:
			if update.ChannelPost != nil {
				chat := update.ChannelPost.Chat
				s.AddChat(chat.ID, chat.UserName, chat.Title)
				break
			}
			if update.EditedMessage != nil {
				if allowed := checkAccess(update.EditedMessage); !allowed {
					accessDeniedMessage(update.EditedMessage.Chat.ID)
					break
				}
				chatID := update.EditedMessage.Chat.ID
				msgID := update.EditedMessage.MessageID
				if tpl, ok := s.GetTemplateBySource(s.MessagePtr{ChatID: chatID, MessageID: msgID}); ok {
					tpl.Text = rich.MessageToHTML(update.EditedMessage)
					if ok := s.SetTemplate(tpl); ok {
						sync = true
					} else {
						msg := t.NewMessage(chatID, "Ошибка в шаблоне.")
						msg.ParseMode = "html"
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
				msg.ParseMode = "html"
				bot.Send(msg)
			}
		}
	}

}

func checkAccess(message *t.Message) bool {
	username := message.From.UserName
	for _, admin := range cfg.Config.Admins {
		if username == admin {
			return true
		}
	}
	return false
}

func accessDeniedMessage(chatID int64) {
	msg := t.NewMessage(chatID, "Я с тобой не разговариваю, обратись к админам.")
	msg.ParseMode = "html"
	bot.Send(msg)
}

func Sync() {
	pushTemplates()
	s.PersistTemplates()
}

func getChatByName(name string) (t.Chat, error) {
	return bot.GetChat(t.ChatConfig{SuperGroupUsername: name})
}

func getChatByID(id int64) (t.Chat, error) {
	return bot.GetChat(t.ChatConfig{ChatID: id})
}

func getChat(addr string) (t.Chat, error) {
	if len(addr) == 0 {
		return t.Chat{}, errors.New("No channel address specified")
	}
	if addr[0] == '@' {
		return getChatByName(addr)
	} else {
		chId, err := strconv.ParseInt(addr, 10, 64)
		if err != nil {
			return t.Chat{}, errors.New("Malformed channel address")
		}
		return getChatByID(chId)
	}
}

func pushTemplates() {
	log.Println("Push templates")
	for _, tpl := range s.GetTemplates() {
		if tpl.IsPosted() {
			if tpl.TargetMessagePtr.ChatID == 0 {
				chat, err := getChat(tpl.TargetChannel)
				if err != nil {
					log.Println(err)
					continue
				}
				tpl.TargetMessagePtr.ChatID = chat.ID
				s.SetTemplate(tpl)
			}
			edit := t.EditMessageTextConfig{
				BaseEdit: t.BaseEdit{
					ChatID:    tpl.TargetMessagePtr.ChatID,
					MessageID: tpl.TargetMessagePtr.MessageID,
				},
				Text: tpl.Apply(srv.GetVariablesValues()),
			}
			edit.ParseMode = "html"
			if _, err := bot.Send(edit); err != nil {
				log.Println("Message edit error:", err)
			}
		} else {
			chat, err := getChat(tpl.TargetChannel)
			if err != nil {
				log.Println(err)
				continue
			}
			msg := t.NewMessage(chat.ID, tpl.Apply(srv.GetVariablesValues()))
			msg.ParseMode = "html"
			postedMsg, err := bot.Send(msg)
			if err == nil {
				tpl.TargetMessagePtr = s.MessagePtr{ChatID: chat.ID, MessageID: postedMsg.MessageID}
				s.SetTemplate(tpl)
			} else {
				log.Println("Message add error:", err)
			}
		}
	}
}

func makeUpdateConfig(timeout int) t.UpdateConfig {
	cfg := t.NewUpdate(0)
	cfg.Timeout = timeout
	return cfg
}
