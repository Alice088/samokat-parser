package dto

type Subcategory struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Slug              string `json:"slug"`
	ProductCategories *[]*ProductCategory
}
