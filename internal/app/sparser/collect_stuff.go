package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	errs "alice088/sparser/internal/pkg/errors"
	"alice088/sparser/internal/pkg/samokat"
	"context"
	"encoding/json"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"os"
	"strconv"
	"strings"
	"time"
)

func (p *Parser) CollectStuff(geo *dto.GEO) (*[]dto.Category, error) {
	var categories *[]dto.Category
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
	defer cancel()

	err = p.getCategories(ctx, geoCookie, categories)

	p.Log.Info().Interface("Categories", categories).Msg("Categories") //todo слишком быстро отдает, нет запроса или горутина не успевает
	os.Exit(1)

	if err != nil {
		return nil, err
	}

	err = chromedp.Run(ctx,
		geoCookie,
		chromedp.Navigate(samokat.MAIN),
		chromedp.Sleep(2*time.Second),
		chromedp.Click(`(//a[contains(@href, '/category/molochnoe-i-yaytsa')])`, chromedp.BySearch),
		chromedp.Sleep(1000*time.Second),
	)
	if err != nil {
		log.Fatal().Err(err).Send()
		return nil, err
	}

	return categories, &errs.ErrSessionDataMissing{}
}

func (p *Parser) getCategories(ctx context.Context, geoCookie *network.SetCookieParams, categories *[]dto.Category) error {
	return chromedp.Run(ctx,
		geoCookie,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.ListenTarget(ctx, func(event interface{}) {
				switch ev := event.(type) {
				case *network.EventResponseReceived:
					if strings.Contains(ev.Response.URL, "list") {
						log.Debug().Msgf("List caought: %s", ev.Response.URL)

						go func() {
							body, err := network.GetResponseBody(ev.RequestID).Do(ctx)

							if err != nil {
								log.Fatal().Err(err).Msgf("Failed to get response body")
							}

							for i := 1; i < 18; i++ {
								jsonCategory := gjson.Get(string(body), strconv.Itoa(i))

								if jsonCategory.Exists() {
									log.Debug().Msgf("No category in index range(1-17): ?")
									continue
								}

								c := &dto.Category{}
								err := json.Unmarshal([]byte(jsonCategory.String()), c)
								if err != nil {
									p.Log.Error().Err(err).Int("INDEX", i).Msgf("Failed to unmarshal category")
									continue
								}
								*categories = append(*categories, *c)
							}
						}()
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
		chromedp.Flag("headless", true),
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)
}
