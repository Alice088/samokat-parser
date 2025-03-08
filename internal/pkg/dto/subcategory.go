package dto

import "github.com/tidwall/gjson"

type Subcategory struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Slug              string `json:"slug"`
	ProductCategories *[]*ProductCategory
}

func HasSubCategoryProperties(json gjson.Result) bool {
	return !json.Get("id").Exists() && !json.Get("name").Exists()
}

func NewSubCategory(json gjson.Result) *Subcategory {
	if !json.Get("slug").Exists() {
		return &Subcategory{
			Id:                json.Get("id").String(),
			Name:              json.Get("name").String(),
			Slug:              json.Get("id").String(),
			ProductCategories: &[]*ProductCategory{},
		}
	}

	return &Subcategory{
		Id:                json.Get("id").String(),
		Name:              json.Get("name").String(),
		Slug:              json.Get("slug").String(),
		ProductCategories: &[]*ProductCategory{},
	}
}
