package state

type Template struct {
	ID              int
	TargetChannel   string
	SourceMessageID int
	TargetMessageID int
	Text            string
}

func (tpl Template) IsPosted() bool {
	return tpl.TargetMessageID > 0
}
