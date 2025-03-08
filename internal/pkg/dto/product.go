package dto

import (
	"fmt"
	"github.com/tidwall/gjson"
)

type Product struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Link    string `json:"link"`
	Current int    `json:"current"`
}

func HasProductProperties(json gjson.Result) bool {
	return !json.Get("name").Exists() && !json.Get("prices.current").Exists() && !json.Get("uuid").Exists()
}

func HasComboProductProperties(json gjson.Result) bool {
	return !json.Get("title").Exists() && !json.Get("priceFrom").Exists() && !json.Get("id").Exists()
}

func NewProduct(json gjson.Result, isCombo bool) *Product {
	if isCombo {
		return &Product{
			Name:    json.Get("title").String(),
			Current: int(json.Get("priceFrom").Int()) / 100,
			UUID:    json.Get("id").String(),
			Link:    fmt.Sprintf("https://samokat.ru/combo/%s", json.Get("id").String()),
		}
	}

	return &Product{
		Name:    json.Get("name").String(),
		Current: int(json.Get("prices.current").Int()) / 100,
		UUID:    json.Get("uuid").String(),
		Link:    fmt.Sprintf("https://samokat.ru/product/%s", json.Get("uuid").String()),
	}
}
