package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"github.com/tidwall/gjson"
	"strconv"
)

func (p *Parser) fillCategories(body []byte, categories *[]*dto.Category) {
	for i := 1; i < 18; i++ {
		json := gjson.Get(string(body), strconv.Itoa(i))

		if !json.Exists() {
			p.Log.Debug().Msgf("No category in index range(1-17): ?")
			continue
		}

		if dto.HasCategoryProperties(json) {
			p.Log.Debug().Int("INDEX", i).Msgf("Category doesn't have target properties")
			continue
		}

		*categories = append(*categories, dto.NewCategory(json))
	}
}
