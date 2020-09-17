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
