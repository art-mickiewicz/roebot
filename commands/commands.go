package commands

import (
	rich "19u4n4/roebot/richtext"
	srv "19u4n4/roebot/services"
	s "19u4n4/roebot/state"
	u "19u4n4/roebot/util"
	"fmt"
	_ "log"
	"strconv"
	"strings"

	"robpike.io/filter"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
	strip "github.com/grokify/html-strip-tags-go"
)

var DefaultTransition = func(m *t.Message) Handler {
	return Zero{Message: m}
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
	case "template", "templates":
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
			if len(args) < 1 {
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
		for _, descr := range srv.GetVariablesInfo() {
			reply += fmt.Sprintln(descr)
		}
		if len(reply) > 0 {
			r = str(fmt.Sprintf("<pre>%s</pre>", reply))
		}
	case "chats":
		r = cmdChats()
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
		if tpl, ok := s.GetTemplateByID(cmd.TemplateID); ok {
			tpl.Text = rich.MessageToHTML(cmd.Message)
			if ok := s.SetTemplate(tpl); ok {
				r = str("Шаблон установлен")
				sync = true
			} else {
				r = str("Ошибка в шаблоне.")
			}
		} else {
			r = str("Шаблона с таким ID не найдено.")
		}
	} else {
		srcPtr := s.MessagePtr{ChatID: cmd.Message.Chat.ID, MessageID: cmd.Message.MessageID}
		tpl := s.NewTemplate(cmd.TargetChannel, srcPtr, rich.MessageToHTML(cmd.Message))
		if cmd.MessageID > 0 {
			tpl.TargetMessagePtr = s.MessagePtr{ChatID: 0, MessageID: cmd.MessageID}
		}
		if ok := s.SetTemplate(tpl); ok {
			r = str("Шаблон установлен")
			sync = true
		} else {
			r = str("Ошибка в шаблоне.")
		}
	}
	return
}

/* Template commands */

func cmdTemplateList() str {
	if s.GetTemplatesCount() == 0 {
		return "Список шаблонов пуст."
	}

	maxChNameLen := len([]rune("Канал"))
	for _, tpl := range s.GetTemplates() {
		l := len([]rune(tpl.PrettyTarget()))
		if l > maxChNameLen {
			maxChNameLen = l
		}
	}

	titleCh := u.PadLine("Канал", maxChNameLen, " ")
	titleLine := u.PadLine("", maxChNameLen, "=")
	msg := fmt.Sprintln("ID  " + titleCh + "  Текст")
	msg += fmt.Sprintln("==  " + titleLine + "  =======================")
	for _, tpl := range s.GetTemplates() {
		text := strip.StripTags(tpl.Text)
		row := fmt.Sprintf("%2d  %s  %s", tpl.ID, u.PadLine(tpl.PrettyTarget(), maxChNameLen, " "), u.TrimLine(text, 20))
		msg += fmt.Sprintln(row)
	}
	return str(fmt.Sprintf("<pre>%s</pre>", msg))
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

func cmdChats() str {
	if s.GetChatsCount() == 0 {
		return "Список чатов пуст."
	}

	maxChUsernameLen := len([]rune("Пользователь"))
	maxChTitleLen := len([]rune("Заголовок"))
	for _, ch := range s.GetChats() {
		ul := len([]rune(ch.Username))
		tl := len([]rune(ch.Title))
		if ul > maxChUsernameLen {
			maxChUsernameLen = ul
		}
		if tl > maxChTitleLen {
			maxChTitleLen = tl
		}
	}

	idCh := u.PadLine("ID", 16, " ")
	idLine := u.PadLine("", 16, "=")
	userCh := u.PadLine("Пользователь", maxChUsernameLen, " ")
	userLine := u.PadLine("", maxChUsernameLen, "=")
	titleCh := u.PadLine("Заголовок", maxChTitleLen, " ")
	titleLine := u.PadLine("", maxChTitleLen, "=")
	msg := fmt.Sprintln(idCh + "  " + userCh + "  " + titleCh)
	msg += fmt.Sprintln(idLine + "  " + userLine + "  " + titleLine)
	for _, ch := range s.GetChats() {
		row := fmt.Sprintf(
			"%10d  %s  %s",
			ch.ID, u.PadLine(ch.Username, maxChUsernameLen, " "),
			ch.Title,
		)
		msg += fmt.Sprintln(row)
	}
	return str(fmt.Sprintf("<pre>%s</pre>", msg))
}
