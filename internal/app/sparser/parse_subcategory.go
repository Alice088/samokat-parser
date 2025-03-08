package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"github.com/tidwall/gjson"
	"strconv"
)

func (p *Parser) parseSubcategory(categories *[]*dto.Category, body []byte) {
	for _, category := range *categories {
		for i := 22; i < 121; i++ {
			json := gjson.Get(string(body), strconv.Itoa(i))
			if !json.Get("parentId").Exists() {
				p.Log.Debug().Msgf("No subcategory in index range(18-120): ?")
				continue
			}

			if dto.HasSubCategoryProperties(json) {
				p.Log.Debug().Int("INDEX", i).Msgf("Subcategory doesn't has target property")
				continue
			}

			if json.Get("parentId").String() == category.Id {
				*category.Subcategories = append(*category.Subcategories, dto.NewSubCategory(json))
			}
		}
	}
}
