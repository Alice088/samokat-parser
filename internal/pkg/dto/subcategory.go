package dto

type Subcategory struct {
	Id                string `json:"id"`
	ParentId          string `json:"parentId"`
	Name              string `json:"name"`
	ProductCategories []*ProductCategory
}
