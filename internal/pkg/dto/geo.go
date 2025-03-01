package dto

type GEO struct {
	ID       int     `json:"id"`
	Region   string  `json:"region"`
	City     string  `json:"city"`
	Street   string  `json:"street"`
	District string  `json:"district"`
	House    string  `json:"house"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}
