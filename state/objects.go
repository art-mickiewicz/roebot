package state

type Template struct {
	ID               int
	TargetChannel    string
	SourceMessagePtr MessagePtr
	TargetMessagePtr MessagePtr
	Text             string
}

type MessagePtr struct {
	ChatID    int64
	MessageID int
}

func (tpl Template) IsPosted() bool {
	return tpl.TargetMessagePtr.MessageID > 0
}
