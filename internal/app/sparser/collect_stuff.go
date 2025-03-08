package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	errs "alice088/sparser/internal/pkg/errors"
	"alice088/sparser/internal/pkg/samokat"
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
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

	ctx, cancel := chromedp.NewContext(allocCtx)
	parsingContext.Ctx = ctx
	defer cancel()

	err = p.getCategories(parsingContext, geoCookie, categories)

	if err != nil {
		return nil, err
	}

	err = chromedp.Run(ctx,
		geoCookie,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(samokat.MAIN),
		chromedp.Sleep(2*time.Second),
		chromedp.Click(`(//a[contains(@href, '/category/molochnoe-i-yaytsa')])`, chromedp.BySearch),
		chromedp.Sleep(1000*time.Second),
	)
	p.Log.Info().Interface("Categories", (*categories)[0]).Msg("Categories")
	if err != nil {
		log.Fatal().Err(err).Send()
		return nil, err
	}

	return categories, &errs.ErrSessionDataMissing{}
}

func (p *Parser) getCategories(parsingContext *dto.ParsingContext, geoCookie *network.SetCookieParams, categories *[]*dto.Category) error {
	return chromedp.Run(parsingContext.Ctx,
		geoCookie,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.ListenTarget(ctx, func(event interface{}) {
				switch ev := event.(type) {
				case *network.EventResponseReceived:
					if strings.Contains(ev.Response.URL, "categories/list") && !(*parsingContext.Skip)["categories/list"] {
						p.Log.Debug().Msgf("List caught: %s", ev.Response.URL)
						p.parseList(ev, parsingContext, categories)
					}
				}

			})
			return nil
		}))
}

func (p *Parser) setupChromedpOptions() []chromedp.ExecAllocatorOption {
	return append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(os.Getenv("EXEC_CHROME_PATH")),
		chromedp.UserAgent(samokat.USER_AGENT),
		chromedp.Flag("headless", false),
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)
}
