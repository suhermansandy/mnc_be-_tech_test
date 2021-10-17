package model

type Payment struct {
	Default
	From    string  `json:"from"`
	To      string  `json:"to"`
	Nominal float64 `json:"nominal"`
}
