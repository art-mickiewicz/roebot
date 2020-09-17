package main

import (
	"encoding/json"
	"log"
	"net/url"

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
	chans, err := sync.getAdminedChannels()
	if err != nil {
		log.Println(err)
		return
	}
	for _, ch := range chans {
		log.Println(ch.UserName)
	}

	// for _, tpl := range sync.templates {

	// }
}

func (sync synchronizer) getAdminedChannels() ([]t.Chat, error) {
	log.Println("Getting admined channels!!!")
	v := url.Values{}
	resp, err := sync.bot.MakeRequest("c1hannels.getAdminedPublicChannels", v)
	if err != nil {
		return make([]t.Chat, 0), err
	}
	log.Println("Try print result")
	log.Println(resp.Result)

	var chats []t.Chat
	err = json.Unmarshal(resp.Result, &chats)

	return chats, err
}

func makeUpdateConfig(timeout int) t.UpdateConfig {
	cfg := t.NewUpdate(0)
	cfg.Timeout = timeout
	return cfg
}

// GetChatMember gets a specific chat member.
// func (bot *BotAPI) GetChatMember(config ChatConfigWithUser) (ChatMember, error) {
// 	v := url.Values{}

// 	if config.SuperGroupUsername == "" {
// 		v.Add("chat_id", strconv.FormatInt(config.ChatID, 10))
// 	} else {
// 		v.Add("chat_id", config.SuperGroupUsername)
// 	}
// 	v.Add("user_id", strconv.Itoa(config.UserID))

// 	resp, err := bot.MakeRequest("getChatMember", v)
// 	if err != nil {
// 		return ChatMember{}, err
// 	}

// 	var member ChatMember
// 	err = json.Unmarshal(resp.Result, &member)

// 	bot.debugLog("getChatMember", v, member)

// 	return member, err
// }
