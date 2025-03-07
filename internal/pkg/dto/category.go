package dto

type Category struct {
	Name          string `json:"name"`
	SubCategories []*Subcategory
}
