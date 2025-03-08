package dto

import "github.com/tidwall/gjson"

type Category struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Subcategories *[]*Subcategory
}

func HasCategoryProperties(json gjson.Result) bool {
	return !json.Get("id").Exists() && !json.Get("name").Exists()
}

func NewCategory(json gjson.Result) *Category {
	return &Category{
		Id:            json.Get("id").String(),
		Name:          json.Get("name").String(),
		Subcategories: &[]*Subcategory{},
	}
}
