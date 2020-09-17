package state

type Template struct {
	ID              int
	TargetChannel   string
	SourceMessageID int
	TargetMessageID int
	Text            string
}
