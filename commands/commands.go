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

type Replier interface {
	Reply(bot *t.BotAPI) string
}

type str string

func (s str) Reply(bot *t.BotAPI) string {
	return string(s)
}

type Transition func(*t.Message) Handler

var DefaultTransition = func(m *t.Message) Handler {
	return Zero{Message: m}
}

type Handler interface {
	Handle() (transitTo Transition, r Replier, sync bool)
}

type Zero struct{ Message *t.Message }

func (cmd Zero) Handle() (transitTo Transition, r Replier, sync bool) {
	transitTo = DefaultTransition
	r = str("")
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
			r = str("Не указана субкоманда: { list | add | edit | delete }")
			return
		}
		subcmd := args[0]
		args = args[1:]
		switch subcmd {
		case "list":
			r = cmdTemplateList()
		case "add":
			if len(args) < 1 || args[0][0] != '@' {
				r = str("Не указан канал для отправки шаблона первым аргументом к команде \"template add\".")
				return
			}
			chName := args[0][1:]
			msgID := 0
			if len(args) > 1 {
				msgID, _ = strconv.Atoi(args[1])
			}
			transitTo, r = cmdTemplateAdd(chName, msgID)
		case "edit":
			if len(args) < 1 {
				r = str("Не указан ID шаблона для редактирования первым аргументом к команде \"template edit\".")
				return
			}
			if chID, err := strconv.Atoi(args[0]); err == nil {
				transitTo, r = cmdTemplateEdit(chID)
			} else {
				r = str("Неверный тип аргумента ID шаблона.")
			}
		case "delete":
			if len(args) < 1 {
				r = str("Не указан ID шаблона для редактирования первым аргументом к команде \"template delete\".")
				return
			}
			if chID, err := strconv.Atoi(args[0]); err == nil {
				r = cmdTemplateDelete(chID)
			} else {
				r = str("Неверный тип аргумента ID шаблона.")
			}
		}
	case "help", "start":
		if len(args) < 1 {
			r = GetHelp("")
		} else {
			help := GetHelp(args[0])
			if help == "" {
				r = str("Невозможно получить справку - команда не найдена.")
			} else {
				r = help
			}
		}
	case "variables":
		var reply string
		for k, _ := range srv.GetVariablesInfo() {
			reply += fmt.Sprintln(fmt.Sprintf("%s", k))
		}
		if len(reply) > 0 {
			r = str(fmt.Sprintln("```") + reply + fmt.Sprintln("```"))
		}
	case "messages":

	}
	return
}

type SetTemplate struct {
	Message       *t.Message
	TargetChannel string
	TemplateID    int
	MessageID     int
}

func (cmd SetTemplate) Handle() (transitTo Transition, r Replier, sync bool) {
	transitTo = DefaultTransition
	sync = false
	if cmd.TemplateID > 0 {
		if tpl := s.GetTemplateByID(cmd.TemplateID); tpl != nil {
			tpl.Text = cmd.Message.Text
			r = str("Шаблон установлен")
			sync = true
		} else {
			r = str("Шаблона с таким ID не найдено.")
		}
	} else {
		srcPtr := s.MessagePtr{ChatID: cmd.Message.Chat.ID, MessageID: cmd.Message.MessageID}
		tpl := s.NewTemplate(cmd.TargetChannel, srcPtr, cmd.Message.Text)
		if cmd.MessageID > 0 {
			tpl.TargetMessagePtr = s.MessagePtr{ChatID: 0, MessageID: cmd.MessageID}
		}
		s.Templates = append(s.Templates, tpl)
		r = str("Шаблон установлен")
		sync = true
	}
	return
}

/* Template commands */

func cmdTemplateList() str {
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
	msg += fmt.Sprintln("==  " + titleLine + "  =======================")
	for _, tpl := range s.Templates {
		row := fmt.Sprintf("%2d  %s  %s", tpl.ID, u.PadLine(tpl.TargetChannel, maxChNameLen, " "), u.TrimLine(tpl.Text, 20))
		msg += fmt.Sprintln(row)
	}
	return str(fmt.Sprintln("```") + msg + fmt.Sprintln("```"))
}

func cmdTemplateAdd(channelName string, msgID int) (transitTo Transition, reply str) {
	transitTo = func(m *t.Message) Handler {
		return SetTemplate{Message: m, TargetChannel: channelName, MessageID: msgID}
	}
	reply = "Жду шаблон объявления следующим сообщением."
	return
}

func cmdTemplateEdit(templateID int) (transitTo Transition, reply str) {
	transitTo = func(m *t.Message) Handler {
		return SetTemplate{Message: m, TemplateID: templateID}
	}
	reply = "Жду шаблон объявления следующим сообщением."
	return
}

func cmdTemplateDelete(templateID int) str {
	deleted := s.DeleteTemplateByID(templateID)
	if deleted > 0 {
		return "Удалено."
	} else {
		return "Шаблона с таким ID не найдено."
	}
}
