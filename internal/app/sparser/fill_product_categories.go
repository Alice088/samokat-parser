package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"github.com/tidwall/gjson"
)

func (p *Parser) fillProductCategory(body []byte, subcategory *dto.Subcategory) {
	json := gjson.Get(string(body), "categories")

	if !json.Exists() {
		p.Log.Error().Msgf("Categories in %s not found", subcategory.Name)
		return
	}

	json.ForEach(func(key gjson.Result, value gjson.Result) bool {
		if !value.Get("name").Exists() {
			p.Log.Debug().Msgf("Product category doesn't have target properties")
			return false
		}

		*((*subcategory).ProductCategories) = append(*(*subcategory).ProductCategories, dto.NewProductCategory(value))
		return true
	})
}
