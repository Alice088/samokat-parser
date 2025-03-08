package dto

import (
	"encoding/json"
	"net/url"
	"strings"
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

func (g GEO) URL() (string, error) {
	bytes, err := json.Marshal(g)

	if err != nil {
		return "", err
	}

	escaped := url.PathEscape(string(bytes))

	escaped = strings.ReplaceAll(escaped, "%7B", "{")
	escaped = strings.ReplaceAll(escaped, "%7D", "}")
	escaped = strings.ReplaceAll(escaped, "%2F", "/")

	return escaped, nil
}
