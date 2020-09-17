package commands

import (
	s "19u4n4/roebot/state"
	"strings"

	"robpike.io/filter"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CommandMode int

const (
	ModeSync        CommandMode = -1
	ModeZero        CommandMode = 0
	ModeSetTemplate CommandMode = 1
)

type Handler interface {
	Handle() (transitTo CommandMode, reply string)
}

type Zero struct{ Message *t.Message }

func (cmd Zero) Handle() (transitTo CommandMode, reply string) {
	transitTo = ModeZero
	reply = ""
	args := filter.Choose(strings.Split(cmd.Message.Text, " "), func(x string) bool {
		return x != ""
	}).([]string)
	name, args := args[0], args[1:]
	switch name {
	case "template":
		transitTo = ModeSetTemplate
		reply = "Жду шаблон объявления следующим сообщением."
	case "join":
		if len(args) != 1 || args[0][0] != '@' {
			reply = "Не указано имя чата для подключения."
			break
		}
		chName := args[0]
		if _, ok := s.Chats[chName]; ok {
			reply = "Кажется, я там уже был."
			transitTo = ModeSync
			break
		}
		reply = "Подключаюсь к " + chName
		s.Chats[chName] = s.Chat{ID: 0, Name: chName}
		transitTo = ModeSync
	}
	return
}

type SetTemplate struct{ Message *t.Message }

func (cmd SetTemplate) Handle() (transitTo CommandMode, reply string) {
	tpl := s.Template{ID: 0, SourceMessageID: 0, TargetMessageID: 0, Text: ""}
	s.Templates = append(s.Templates, tpl)
	transitTo = ModeSync
	reply = "Шаблон установлен"
	return
}
