package richtext

import (
	"19u4n4/roebot/richtext/style"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Entity struct {
	Style style.Style
	Text  string
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

func EntitiesToHTML([]Entity) string {
	return ""
}
