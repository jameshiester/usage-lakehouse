package model

type Address struct {
	Attention    *string `json:"attention"`
	AddressLine1 string  `json:"address_line_1"`
	AddressLine2 *string `json:"address_line_2"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	Zip          string  `json:"zip"`
	Country      string  `json:"country"`
}
