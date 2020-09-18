package commands

import (
	srv "19u4n4/roebot/services"
	s "19u4n4/roebot/state"
	u "19u4n4/roebot/util"
	"fmt"
	_ "log"
	"strconv"
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
		if len(args) < 1 {
			reply = "Не указана субкоманда: { list | add | edit | delete }"
			return
		}
		subcmd := args[0]
		args = args[1:]
		switch subcmd {
		case "list":
			reply = cmdTemplateList()
		case "add":
			if len(args) < 1 || args[0][0] != '@' {
				reply = "Не указан канал для отправки шаблона первым аргументом к команде \"template add\"."
				return
			}
			transitTo, reply = cmdTemplateAdd(args[0][1:])
		case "edit":
			if len(args) < 1 || args[0][0] != '@' {
				reply = "Не указан ID шаблона для редактирования первым аргументом к команде \"template edit\"."
				return
			}
			if chID, ok := strconv.Atoi(args[0][1:]); ok != nil {
				transitTo, reply = cmdTemplateEdit(chID)
			} else {
				reply = "Неверный тип аргумента ID шаблона."
				return
			}
		case "delete":
			if chID, ok := strconv.Atoi(args[0][1:]); ok != nil {
				reply = cmdTemplateDelete(chID)
			} else {
				reply = "Неверный тип аргумента ID шаблона."
				return
			}
		}
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

/* Template commands */

func cmdTemplateList() string {
	if len(s.Templates) == 0 {
		return "Список шаблонов пуст."
	}

	maxChNameLen := 0
	for _, tpl := range s.Templates {
		l := len(tpl.TargetChannel)
		if l > maxChNameLen {
			maxChNameLen = l
		}
	}

	titleCh := u.PadLine("Канал", maxChNameLen, " ")
	titleLine := u.PadLine("", maxChNameLen, "=")
	msg := fmt.Sprintln("ID  " + titleCh + "  Текст")
	msg += fmt.Sprintln("==  " + titleLine + "  =============")
	for _, tpl := range s.Templates {
		row := fmt.Sprintf("%2d  %s  %s", tpl.ID, u.PadLine(tpl.TargetChannel, maxChNameLen, " "), u.TrimLine(tpl.Text, 10))
		msg += fmt.Sprintln(row)
	}
	return fmt.Sprintln("```") + msg + fmt.Sprintln("```")
}

func cmdTemplateAdd(channelName string) (transitTo Transition, reply string) {
	transitTo = func(m *t.Message) Handler {
		return SetTemplate{Message: m, TargetChannel: channelName}
	}
	reply = "Жду шаблон объявления следующим сообщением."
	return
}

func cmdTemplateEdit(templateID int) (transitTo Transition, reply string) {
	transitTo = func(m *t.Message) Handler {
		return SetTemplate{Message: m, TemplateID: templateID}
	}
	reply = "Жду шаблон объявления следующим сообщением."
	return
}

func cmdTemplateDelete(templateID int) string {
	return "Удалено."
}
