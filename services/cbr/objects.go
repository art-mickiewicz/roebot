package cbr

type CBRResponse struct {
	Date         string
	PreviousDate string
	PreviousURL  string
	Timestamp    string
	Valute       map[string]Symbol
}

type Symbol struct {
	ID       string
	NumCode  string
	CharCode string
	Nominal  int
	Name     string
	Value    float64
	Previous float64
}
