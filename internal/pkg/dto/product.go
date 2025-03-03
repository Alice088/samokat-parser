package dto

type Product struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Prices prices `json:"prices"`
}

type prices struct {
	Current int `json:"current"`
}
