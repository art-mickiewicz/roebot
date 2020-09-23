package richtext

import (
	"19u4n4/roebot/richtext/style"
	"fmt"
	"unicode/utf16"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Token struct {
	Style     style.Style
	Text      string
	Subtokens []Token
}

func (tok Token) String() string {
	switch tok.Style {
	case style.Plain:
		return tok.Text
	case style.Bold:
		return fmt.Sprintf("<b>%s</b>", tok.Text)
	case style.Italic:
		return fmt.Sprintf("<i>%s</i>", tok.Text)
	case style.Underline:
		return fmt.Sprintf("<u>%s</u>", tok.Text)
	case style.Strikethrough:
		return fmt.Sprintf("<s>%s</s>", tok.Text)
	case style.Code:
		return fmt.Sprintf("<code>%s</code>", tok.Text)
	case style.Pre:
		return fmt.Sprintf("<pre>%s</pre>", tok.Text)
	default:
		return tok.Text
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func messageToTokens(msg *t.Message, index int, lowerBound int, upperBound int) []Token {
	meLen := 0
	if msg.Entities != nil {
		meLen = len(*msg.Entities)
	}
	if meLen == 0 {
		ret := make([]Token, 1)
		ret[0] = Token{Style: style.Plain, Text: msg.Text}
		return ret
	}

	u16s := utf16.Encode([]rune(msg.Text))
	u16len := len(u16s)
	if upperBound < 0 {
		upperBound = u16len
	}
	ret := make([]Token, 0, 2*meLen+1)
	prevCursor := lowerBound
	cursor := lowerBound
	skipSubtokens := false
	fmt.Println(index, "LOWER", lowerBound, "UPPER", upperBound)
	for i, me := range (*msg.Entities)[index:] {
		fmt.Println("ENTITY", me.Type, me.Offset, me.Length)
		fmt.Println(index, "CURSOR", cursor, "OFFSET", me.Offset)
		if cursor >= upperBound {
			return ret
		}

		subtokens := make([]Token, 0, 5)
		if cursor > me.Offset {
			// trigger subtokens
			if skipSubtokens {
				fmt.Println(index, "SKIP SUBTOKEN")
				continue
			} else {
				fmt.Println(index, "--- SCAN SUBTOKENS ---")
				subtokens = messageToTokens(msg, i, prevCursor, cursor)
				skipSubtokens = true
			}
		} else {
			skipSubtokens = false
		}

		if cursor < me.Offset {
			upTo := min(me.Offset, upperBound)
			ent := Token{
				Style: style.Plain,
				Text:  string(utf16.Decode(u16s[cursor:upTo])),
			}
			ret = append(ret, ent)
		}
		upTo := min(me.Offset+me.Length, upperBound)
		ent := Token{
			Style:     style.FromType(me.Type),
			Text:      string(utf16.Decode(u16s[me.Offset:upTo])),
			Subtokens: subtokens,
		}
		ret = append(ret, ent)
		prevCursor = cursor
		cursor = me.Offset + me.Length
	}
	if cursor < upperBound {
		ent := Token{Style: style.Plain, Text: string(utf16.Decode(u16s[cursor:upperBound]))}
		ret = append(ret, ent)
	}
	fmt.Println(ret)
	return ret
}

func TokensToHTML(toks []Token) string {
	ret := ""
	for _, tok := range toks {
		if len(tok.Subtokens) > 0 {
			ret += TokensToHTML(tok.Subtokens)
		} else {
			ret += tok.String()
		}
	}
	return ret
}

func MessageToHTML(msg *t.Message) string {
	return TokensToHTML(messageToTokens(msg, 0, 0, -1))
}
