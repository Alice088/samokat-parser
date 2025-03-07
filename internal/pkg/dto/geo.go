package dto

import (
	"encoding/json"
	"net/url"
)

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

func (g GEO) ToCookie() (string, error) {
	json, err := json.Marshal(g)

	if err != nil {
		return "", err
	}

	return url.QueryEscape(string(json)), nil
}
