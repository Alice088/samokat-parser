package dto

import "github.com/tidwall/gjson"

type Product struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Current int    `json:"current"`
}

func HasProductProperties(json gjson.Result) bool {
	return !json.Get("name").Exists() && !json.Get("prices.current").Exists() && !json.Get("uuid").Exists()
}

func NewProduct(json gjson.Result) *Product {
	return &Product{
		Name:    json.Get("name").String(),
		Current: int(json.Get("prices.current").Int()) / 100,
		UUID:    json.Get("uuid").String(),
	}
}
