package dto

import "github.com/tidwall/gjson"

type ProductCategory struct {
	Name     string `json:"name"`
	Products *[]*Product
}

func NewProductCategory(json gjson.Result) *ProductCategory {
	return &ProductCategory{
		Name:     json.Get("name").String(),
		Products: &[]*Product{},
	}
}
