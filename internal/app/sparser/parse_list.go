package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"github.com/chromedp/cdproto/network"
	"time"
)

func (p *Parser) parseList(ev *network.EventResponseReceived, parsingContext *dto.ParsingContext, categories *[]*dto.Category) {
	time.Sleep(1 * time.Second)
	body, err := network.GetResponseBody(ev.RequestID).Do(parsingContext.Ctx)

	if err != nil {
		p.Log.Debug().Msgf("Skipping list")
		return
	}
	(*parsingContext.Skip)["categories/list"] = true

	p.parseCategory(body, categories)
	p.parseSubcategory(categories, body)
}
