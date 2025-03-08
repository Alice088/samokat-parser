package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	errs "alice088/sparser/internal/pkg/errors"
	"sync"
)

func (p *Parser) CollectStuff(geo *dto.GEO) (*[]*dto.Category, error) {
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
		if i != 5 {
			continue
		}

		for j, subcategory := range *category.Subcategories {
			if j != 2 {
				continue
			}

			wg.Add(1)
			go func() {
				collectProductCtx, collectProductCtxCancel := dto.NewParsingContext(skip, urlCookie)
				defer collectProductCtxCancel()
				p.getProducts(collectProductCtx, subcategory, wg)
			}()
		}
		wg.Wait()
	}

	return categories, &errs.ErrSessionDataMissing{}
}
