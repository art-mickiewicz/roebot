package commands

import "fmt"

var Help = map[string]string{
	"help [COMMAND]": "help [COMMAND] - список команд",
	"template":       "template { list | add | edit | delete | show } - работа с шаблонами",
	"variables":      "variables - список переменных для шаблона",
	"chats":          "chats - список известных чатов",
}

var FullHelp = map[string]string{
	"template": `
		list - список шаблонов
		add CHANNEL [MSGID] - добавить шаблон для канала CHANNEL и сообщения MSGID (или нового сообщения)
		edit TEMPLATE_ID - ввести новую версию шаблона
		delete TEMPLATE_ID - удалить шаблон
		show TEMPLATE_ID - отправить текст шаблона в этот чат
	`,
}

func GetHelp(cmd string) str {
	var help string
	if cmd == "" {
		for _, h := range Help {
			help += fmt.Sprintln(h)
		}
	} else if cmdHelp, ok := Help[cmd]; ok {
		help = cmdHelp
		if fullHelp, ok := FullHelp[cmd]; ok {
			help += fmt.Sprintln(fullHelp)
		}
	}
	return str(fmt.Sprintf("<pre>%s</pre>", help))
}
