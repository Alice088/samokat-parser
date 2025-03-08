package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"github.com/chromedp/cdproto/network"
	"github.com/rs/zerolog/log"
	"time"
)

func (p *Parser) parseProducts(ev *network.EventResponseReceived, parsingContext *dto.ParsingContext, subcategory *dto.Subcategory) {
	time.Sleep(1 * time.Second)
	body, err := network.GetResponseBody(ev.RequestID).Do(parsingContext.EventCtx)

	if len(*subcategory.ProductCategories) != 0 {
		return
	}

	if err != nil {
		if err.Error() == "No resource with given identifier found (-32000)" {
			return
		}

		log.Error().Err(err).Msg("Failed to get response body")
		return
	}

	(*parsingContext.Skip).Store(subcategory.Id, true)

	p.fillProductCategory(body, subcategory)
	p.fillProducts(body, subcategory)

	p.Log.Info().Msg("Finished filling products")

	return
}
