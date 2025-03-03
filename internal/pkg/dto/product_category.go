package dto

type ProductCategory struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Products []*Product
}
