package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"fmt"
	"github.com/tidwall/gjson"
)

func (p *Parser) fillProducts(body []byte, subcategory *dto.Subcategory) {
	for _, productCategory := range *subcategory.ProductCategories {
		json := gjson.Get(string(body), fmt.Sprintf("categories.#(name=%s).comboSets", productCategory.Name))
		if json.Exists() {
			for i, product := range json.Array() {
				if dto.HasComboProductProperties(product) {
					p.Log.Debug().Int("INDEX", i).Msgf("Combo product doesn't have target properties")
					continue
				}
				*productCategory.Products = append(*productCategory.Products, dto.NewProduct(product, true))
			}
			return
		}

		json = gjson.Get(string(body), fmt.Sprintf("categories.#(name=%s).products", productCategory.Name))

		if !json.Exists() {
			p.Log.Error().Msgf("product category %s not found", productCategory.Name)
			return
		}

		for i, product := range json.Array() {
			if dto.HasProductProperties(product) {
				p.Log.Debug().Int("INDEX", i).Msgf("Product doesn't have target properties")
				continue
			}

			*productCategory.Products = append(*productCategory.Products, dto.NewProduct(product, false))
		}
	}
}
