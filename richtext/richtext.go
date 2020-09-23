package richtext

/** Supported HTML subset
 *
 * <b>bold</b>, <strong>bold</strong>
 * <i>italic</i>, <em>italic</em>
 * <u>underline</u>, <ins>underline</ins>
 * <s>strikethrough</s>, <strike>strikethrough</strike>, <del>strikethrough</del>
 * <b>bold <i>italic bold <s>italic bold strikethrough</s> <u>underline italic bold</u></i> bold</b>
 * <a href="http://www.example.com/">inline URL</a>
 * <a href="tg://user?id=123456789">inline mention of a user</a>
 * <code>inline fixed-width code</code>
 * <pre>pre-formatted fixed-width code block</pre>
 * <pre><code class="language-python">pre-formatted fixed-width code block written in the Python programming language</code></pre>
 */

import (
	"19u4n4/roebot/richtext/style"
	"fmt"
	"unicode/utf16"

	t "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Token struct {
	Style     style.Style
	Text      string
	URL       string
	Subtokens []Token
}

func (tok Token) tagWrap(text string) string {
	switch tok.Style {
	case style.Plain:
		return tok.Text
	case style.Bold:
		return fmt.Sprintf("<b>%s</b>", text)
	case style.Italic:
		return fmt.Sprintf("<i>%s</i>", text)
	case style.Underline:
		return fmt.Sprintf("<u>%s</u>", text)
	case style.Strikethrough:
		return fmt.Sprintf("<s>%s</s>", text)
	case style.Code:
		return fmt.Sprintf("<code>%s</code>", text)
	case style.Pre:
		return fmt.Sprintf("<pre>%s</pre>", text)
	case style.Link:
		return fmt.Sprintf("<a href=\"%s\">%s</a>", tok.URL, text)
	default:
		return text
	}
}

func (tok Token) String() string {
	return tok.tagWrap(tok.Text)
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
	tokCount := 0
	// fmt.Println(index, "LOWER", lowerBound, "UPPER", upperBound)
	for i, me := range (*msg.Entities)[index:] {
		// fmt.Println("ENTITY", me.Type, me.Offset, me.Length, "CURSOR", cursor)
		if cursor >= upperBound {
			return ret
		}

		/* Trigger subtokens */
		if cursor > me.Offset {
			if !skipSubtokens {
				// fmt.Println("@ SUBTOKENS @")
				ret[tokCount-1].Subtokens = messageToTokens(msg, i, prevCursor, cursor)
				skipSubtokens = true

			}
			continue
		} else {
			skipSubtokens = false
		}

		/* Add plain token before entity */
		if cursor < me.Offset {
			upTo := min(me.Offset, upperBound)
			tok := Token{
				Style: style.Plain,
				Text:  string(utf16.Decode(u16s[cursor:upTo])),
			}
			ret = append(ret, tok)
			// fmt.Println("@ TOKEN PLAIN:", tok.Text)
			tokCount++
			cursor = me.Offset
		}

		/* Token from entity */
		upTo := min(me.Offset+me.Length, upperBound)
		tok := Token{
			Style: style.FromType(me.Type),
			Text:  string(utf16.Decode(u16s[cursor:upTo])),
			URL:   me.URL,
		}
		ret = append(ret, tok)
		// fmt.Println("@ TOKEN "+me.Type+":", tok.Text)
		tokCount++

		prevCursor = cursor
		cursor = me.Offset + me.Length
	}
	if cursor < upperBound {
		tok := Token{Style: style.Plain, Text: string(utf16.Decode(u16s[cursor:upperBound]))}
		ret = append(ret, tok)
		tokCount++
	}
	// fmt.Println("RETURN", index, ret)
	return ret
}

func TokensToHTML(toks []Token) string {
	ret := ""
	for _, tok := range toks {
		if len(tok.Subtokens) > 0 {
			ret += tok.tagWrap(TokensToHTML(tok.Subtokens))
		} else {
			ret += tok.String()
		}
	}
	return ret
}

func MessageToHTML(msg *t.Message) string {
	return TokensToHTML(messageToTokens(msg, 0, 0, -1))
}
