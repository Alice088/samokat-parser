package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	errs "alice088/sparser/internal/pkg/errors"
	"alice088/sparser/internal/pkg/samokat"
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"os"
	"strings"
	"time"
)

func (p *Parser) CollectStuff(geo *dto.GEO) (*[]*dto.Category, error) {
	categories := &[]*dto.Category{}
	parsingContext := &dto.ParsingContext{
		Skip: &map[string]bool{
			"categories/list": false,
		},
	}
	opts := p.setupChromedpOptions()
	cookieExpr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
	urlCookie, err := geo.ToCookie()
	geoCookie := network.SetCookie("SELECTED_ADDRESS_KEY", urlCookie).
		WithExpires(&cookieExpr).
		WithDomain(samokat.DOMAIN).
		WithPath("/")

	if err != nil {
		p.Log.Fatal().Err(err).Int("GEO_ID", geo.ID).Msgf("Failed to cast geo to cookie")
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, chromeCancel := chromedp.NewContext(allocCtx)
	parsingContext.ChromeCtx = ctx
	defer chromeCancel()

	err = p.getCategories(parsingContext, geoCookie, categories)
	if err != nil {
		p.Log.Fatal().Err(err).Send()
		return nil, err
	}

	p.Log.Info().Interface("Categories", categories).Msg("Categories")
	return categories, &errs.ErrSessionDataMissing{}
}

func (p *Parser) getCategories(parsingContext *dto.ParsingContext, geoCookie *network.SetCookieParams, categories *[]*dto.Category) error {
	return chromedp.Run(parsingContext.ChromeCtx,
		geoCookie,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.ListenTarget(ctx, func(event interface{}) {
				parsingContext.EventCtx = ctx
				switch ev := event.(type) {
				case *network.EventResponseReceived:
					if strings.Contains(ev.Response.URL, "categories/list") && !(*parsingContext.Skip)["categories/list"] {
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
}

func (p *Parser) setupChromedpOptions() []chromedp.ExecAllocatorOption {
	return append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(os.Getenv("EXEC_CHROME_PATH")),
		chromedp.UserAgent(samokat.USER_AGENT),
		chromedp.Flag("headless", true),
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)
}
