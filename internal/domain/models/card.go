package models

type Card struct {
	ID     uint64
	Number string
	CVV    string
	Month  string
	Year   string
	Info   string
}
