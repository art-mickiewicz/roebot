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

	commandMode := cmd.ModeZero
	for {
		select {
		case update := <-ch:
			chatID := update.Message.Chat.ID
			var hdl cmd.Handler

			// username := update.Message.From.UserName
			// text := update.Message.Text
			// log.Printf("[%s] %d %s", username, chatID, text)

			switch commandMode {
			case cmd.ModeSync:
				sync := synchronizer{bot: bot, templates: s.Templates, chats: s.Chats}
				sync.start()
				commandMode = cmd.ModeZero
				continue
			case cmd.ModeZero:
				hdl = cmd.Zero{Message: update.Message}
			case cmd.ModeSetTemplate:
				hdl = cmd.SetTemplate{Message: update.Message}
			}

			newMode, reply := hdl.Handle()
			commandMode = newMode
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
	chats     map[string]s.Chat
}

func (sync synchronizer) start() {
	sync.pushTemplates()
}

func (sync synchronizer) pushTemplates() {
	// for _, tpl := range sync.templates {

	// }
}

func makeUpdateConfig(timeout int) t.UpdateConfig {
	cfg := t.NewUpdate(0)
	cfg.Timeout = timeout
	return cfg
}
