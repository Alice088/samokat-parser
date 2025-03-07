package dto

type Category struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Subcategories *[]*Subcategory
}
