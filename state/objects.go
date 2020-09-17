package state

type Template struct {
	ID              int
	SourceMessageID int
	TargetMessageID int
	Text            string
}

type Chat struct {
	ID   int
	Name string
}
