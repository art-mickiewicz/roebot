package richtext

import (
	"19u4n4/roebot/richtext/style"
	"fmt"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Entity struct {
	Style style.Style
	Text  string
}

func (e Entity) String() string {
	switch e.Style {
	case style.Plain:
		return e.Text
	case style.Bold:
		return fmt.Sprintf("<b>%s</b>", e.Text)
	case style.Italic:
		return fmt.Sprintf("<i>%s</t>", e.Text)
	case style.Underline:
		return fmt.Sprintf("<u>%s</u>", e.Text)
	case style.Strikethrough:
		return fmt.Sprintf("<s>%s</s>", e.Text)
	case style.Code:
		return fmt.Sprintf("<code>%s</code>", e.Text)
	case style.Pre:
		return fmt.Sprintf("<pre>%s</pre>", e.Text)
	default:
		return e.Text
	}
}

func MessageToEntities(msg *t.Message) []Entity {
	runes := []rune(msg.Text)
	meLen := len(*msg.Entities)
	ret := make([]Entity, 0, 2*meLen+1)
	cursor := 0
	for _, me := range *msg.Entities {
		if cursor < me.Offset {
			ent := Entity{Style: style.Plain, Text: string(runes[cursor:me.Offset])}
			ret = append(ret, ent)
		}
		ent := Entity{
			Style: style.FromType(me.Type),
			Text:  string(runes[me.Offset : me.Offset+me.Length]),
		}
		ret = append(ret, ent)
		cursor = me.Offset + me.Length
	}
	runelen := len(runes)
	if cursor < runelen {
		ent := Entity{Style: style.Plain, Text: string(runes[cursor:runelen])}
		ret = append(ret, ent)
	}
	return ret
}

func EntitiesToHTML(es []Entity) string {
	ret := ""
	for _, ent := range es {
		ret += ent.String()
	}
	return ret
}
