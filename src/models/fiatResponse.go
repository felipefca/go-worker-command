package models

type FiatResponse struct {
	Code   string `json:"code"`
	CodeIn string `json:"codein"`
	Name   string `json:"name"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Value  string `json:"bid"`
}
