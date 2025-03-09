package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"alice088/sparser/internal/pkg/env"
	"math/rand"
	"sync"
	"time"
)

func (p *Parser) CollectStuff(geo *dto.GEO) *[]*dto.Category {
	rand.Seed(time.Now().UnixNano())

	skip := &sync.Map{}
	skip.Store("categories/list", false)

	urlCookie, err := geo.URL()

	categories := &[]*dto.Category{}
	collectCategoriesCtx, collectCategoriesCtxCancel := dto.NewParsingContext(skip, urlCookie)
	defer collectCategoriesCtxCancel()

	if err != nil {
		p.Log.Fatal().Err(err).Int("GEO_ID", geo.ID).Msgf("Failed to cast geo to cookie")
	}

	p.getCategories(collectCategoriesCtx, categories)

	wg := &sync.WaitGroup{}
	for i, category := range *categories {
		if i >= env.GetCategoryParseLimit() {
			continue
		}

		for _, subcategory := range *category.Subcategories {
			time.Sleep(time.Duration(rand.Intn(10)+5) * time.Second)
			wg.Add(1)
			go func(subcategory dto.Subcategory) {
				defer wg.Done()
				collectProductCtx, collectProductCtxCancel := dto.NewParsingContext(skip, urlCookie)
				defer collectProductCtxCancel()
				p.getProducts(collectProductCtx, &subcategory)
			}(*subcategory)
		}
	}
	wg.Wait()

	return categories
}
