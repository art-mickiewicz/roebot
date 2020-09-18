package commands

import (
	srv "19u4n4/roebot/services"
	s "19u4n4/roebot/state"
	"fmt"
	_ "log"
	"strings"

	"robpike.io/filter"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Transition func(*t.Message) Handler

var DefaultTransition = func(m *t.Message) Handler {
	return Zero{Message: m}
}

type Handler interface {
	Handle() (transitTo Transition, reply string, sync bool)
}

type Zero struct{ Message *t.Message }

func (cmd Zero) Handle() (transitTo Transition, reply string, sync bool) {
	transitTo = DefaultTransition
	reply = ""
	sync = false
	args := filter.Choose(strings.Split(cmd.Message.Text, " "), func(x string) bool {
		return x != ""
	}).([]string)
	name, args := args[0], args[1:]
	if name[0] == '/' {
		name = name[1:]
	}
	switch name {
	case "template":
		if len(args) < 1 || args[0][0] != '@' {
			reply = "Не указан канал для отправки шаблона первым аргументом к команде."
			return
		}
		chName := args[0][1:]
		transitTo = func(m *t.Message) Handler {
			return SetTemplate{Message: m, TargetChannel: chName}
		}
		reply = "Жду шаблон объявления следующим сообщением."
	case "help", "start":
		if len(args) < 1 {
			reply = GetHelp("")
		} else {
			help := GetHelp(args[0])
			if help == "" {
				reply = "Невозможно получить справку - команда не найдена."
			} else {
				reply = help
			}
		}
	case "variables":
		for k, _ := range srv.GetVariablesInfo() {
			reply += fmt.Sprintln(fmt.Sprintf("%s", k))
		}
		if len(reply) > 0 {
			reply = fmt.Sprintln("```") + reply + fmt.Sprintln("```")
		}
	}
	return
}

type SetTemplate struct {
	Message       *t.Message
	TargetChannel string
	TemplateID    int
}

func (cmd SetTemplate) Handle() (transitTo Transition, reply string, sync bool) {
	srcPtr := s.MessagePtr{ChatID: cmd.Message.Chat.ID, MessageID: cmd.Message.MessageID}
	tpl := s.NewTemplate(cmd.TargetChannel, srcPtr, cmd.Message.Text)
	s.Templates = append(s.Templates, tpl)
	transitTo = DefaultTransition
	reply = "Шаблон установлен"
	sync = true
	return
}
