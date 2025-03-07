package dto

type Subcategory struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Slag              string `json:"slag"`
	ProductCategories []*ProductCategory
}
