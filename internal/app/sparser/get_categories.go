package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"alice088/sparser/internal/pkg/samokat"
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"strings"
	"time"
)

func (p *Parser) getCategories(parsingContext *dto.ParsingContext, categories *[]*dto.Category) {
	err := chromedp.Run(parsingContext.ChromeCtx,
		parsingContext.GeoCookie,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.ListenTarget(ctx, func(event interface{}) {
				parsingContext.EventCtx = ctx
				switch ev := event.(type) {
				case *network.EventResponseReceived:
					skip, ok := (*parsingContext.Skip).Load("categories/list")
					if !ok {
						p.Log.Fatal().Msg("Error during get categories/list in map")
					}

					if strings.Contains(ev.Response.URL, "categories/list") && !skip.(bool) {
						p.Log.Debug().Msgf("List caught: %s", ev.Response.URL)
						go p.parseList(ev, parsingContext, categories)
					}
				}

			})
			return nil
		}),
		chromedp.Navigate(samokat.MAIN),
		chromedp.Sleep(4*time.Second),
	)

	if err != nil {
		p.Log.Fatal().Err(err).Msg("Error during parse categories")
	}
}
