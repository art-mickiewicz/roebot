package richtext

import (
	"19u4n4/roebot/richtext/style"
	"fmt"
	"unicode/utf16"

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

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MessageToEntities(msg *t.Message) []Entity {
	u16s := utf16.Encode([]rune(msg.Text))
	u16len := len(u16s)
	meLen := len(*msg.Entities)
	ret := make([]Entity, 0, 2*meLen+1)
	cursor := 0
	for _, me := range *msg.Entities {
		fmt.Println(me.Offset, me.Length)
		if cursor < me.Offset {
			upTo := min(me.Offset, u16len)
			ent := Entity{
				Style: style.Plain,
				Text:  string(utf16.Decode(u16s[cursor:upTo])),
			}
			ret = append(ret, ent)
		}
		upTo := min(me.Offset+me.Length, u16len)
		ent := Entity{
			Style: style.FromType(me.Type),
			Text:  string(utf16.Decode(u16s[me.Offset:upTo])),
		}
		ret = append(ret, ent)
		cursor = me.Offset + me.Length
	}
	if cursor < u16len {
		ent := Entity{Style: style.Plain, Text: string(utf16.Decode(u16s[cursor:u16len]))}
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

func MessageToHTML(msg *t.Message) string {
	return EntitiesToHTML(MessageToEntities(msg))
}
