package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"github.com/chromedp/cdproto/network"
	"github.com/rs/zerolog/log"
	"time"
)

func (p *Parser) parseList(ev *network.EventResponseReceived, parsingContext *dto.ParsingContext, categories *[]*dto.Category) {
	time.Sleep(2 * time.Second)
	body, err := network.GetResponseBody(ev.RequestID).Do(parsingContext.EventCtx)

	if err != nil {
		if err.Error() == "No resource with given identifier found (-32000)" {
			return
		}

		log.Error().Err(err).Msg("Failed to get response body")
		return
	}

	(*parsingContext.Skip).Store("categories/list", true)

	p.parseCategory(body, categories)
	p.parseSubcategory(categories, body)

	return
}
