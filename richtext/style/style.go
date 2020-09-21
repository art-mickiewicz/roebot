package style

type Style int

const (
	Plain         Style = 0
	Bold          Style = 10
	Italic        Style = 20
	Underline     Style = 30
	Strikethrough Style = 40
	Code          Style = 50
	Pre           Style = 60
)

func FromType(t string) Style {
	s, ok := typeMap[t]
	if ok {
		return s
	} else {
		return Plain
	}
}

var typeMap = map[string]Style{
	"bold":          Bold,
	"italic":        Italic,
	"underline":     Underline,
	"strikethrough": Strikethrough,
	"code":          Code,
	"pre":           Pre,
}
