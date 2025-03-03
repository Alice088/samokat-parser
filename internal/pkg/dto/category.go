package dto

type Category struct {
	Id            string `json:"id"`
	ParentId      string `json:"parentId"`
	Name          string `json:"name"`
	SubCategories []*Subcategory
}
