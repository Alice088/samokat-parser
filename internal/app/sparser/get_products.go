package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"alice088/sparser/internal/pkg/samokat"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"strings"
	"sync"
	"time"
)

func (p *Parser) getProducts(parsingContext *dto.ParsingContext, subcategory *dto.Subcategory, wg *sync.WaitGroup) {
	(*parsingContext.Skip).Store(subcategory.Id, false)
	defer wg.Done()

	err := chromedp.Run(parsingContext.ChromeCtx,
		parsingContext.GeoCookie,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.ListenTarget(ctx, func(event interface{}) {
				parsingContext.EventCtx = ctx
				switch ev := event.(type) {
				case *network.EventResponseReceived:
					skip, ok := (*parsingContext.Skip).Load(subcategory.Id)
					if !ok {
						p.Log.Fatal().Msgf("Error during get %s in map", subcategory.Id)
					}

					if strings.Contains(ev.Response.URL, "/categories/"+subcategory.Id) && !skip.(bool) {
						p.Log.Debug().Msg(ev.Response.URL)
						go p.parseProducts(ev, parsingContext, subcategory)
					}
				}

			})
			return nil
		}),
		chromedp.Navigate(samokat.MAIN),
		chromedp.Sleep(2*time.Second),
		chromedp.Reload(),
		chromedp.Sleep(2*time.Second),
		chromedp.Click(fmt.Sprintf(`(//a[contains(@href, '/category/%s')])`, subcategory.Slug), chromedp.BySearch),
		chromedp.Sleep(8*time.Second),
	)

	if err != nil {
		p.Log.Error().Err(err).Msg("Error during parse subcategories")
	}
}
