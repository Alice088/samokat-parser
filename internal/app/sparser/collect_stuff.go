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
	"github.com/tidwall/gjson"
	"os"
	"strconv"
	"strings"
	"time"
)

func (p *Parser) CollectStuff(geo *dto.GEO) (*[]*dto.Category, error) {
	categories := &[]*dto.Category{}
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

	if err != nil {
		return nil, err
	}

	err = chromedp.Run(ctx,
		geoCookie,
		chromedp.Navigate(samokat.MAIN),
		chromedp.Sleep(2*time.Second),
		chromedp.Click(`(//a[contains(@href, '/category/molochnoe-i-yaytsa')])`, chromedp.BySearch),
		chromedp.Sleep(10*time.Second),
	)
	p.Log.Info().Interface("Categories", (*categories)[0]).Msg("Categories")
	if err != nil {
		log.Fatal().Err(err).Send()
		return nil, err
	}

	return categories, &errs.ErrSessionDataMissing{}
}

func (p *Parser) getCategories(ctx context.Context, geoCookie *network.SetCookieParams, categories *[]*dto.Category) error {
	skip := map[string]bool{
		"categories/list": false,
	}

	return chromedp.Run(ctx,
		geoCookie,
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.ListenTarget(ctx, func(event interface{}) {
				switch ev := event.(type) {
				case *network.EventResponseReceived:
					if strings.Contains(ev.Response.URL, "categories/list") && !skip["categories/list"] {
						p.Log.Debug().Msgf("List caught: %s", ev.Response.URL)

						go func() {
							time.Sleep(1 * time.Second)
							body, err := network.GetResponseBody(ev.RequestID).Do(ctx)

							if err != nil {
								p.Log.Debug().Msgf("Skipping list")
								return
							}
							skip["categories/list"] = true

							for i := 1; i < 18; i++ {
								jsonCategory := gjson.Get(string(body), strconv.Itoa(i))

								if !jsonCategory.Exists() {
									p.Log.Debug().Msgf("No category in index range(1-17): ?")
									continue
								}

								if !jsonCategory.Get("id").Exists() && !jsonCategory.Get("name").Exists() {
									p.Log.Debug().Int("INDEX", i).Msgf("Category doesn't has target property")
									continue
								}

								*categories = append(*categories, &dto.Category{
									Id:            jsonCategory.Get("id").String(),
									Name:          jsonCategory.Get("name").String(),
									Subcategories: &[]*dto.Subcategory{},
								})
							}

							for _, category := range *categories {
								for i := 22; i < 121; i++ {
									jsonSubCategory := gjson.Get(string(body), strconv.Itoa(i))
									if !jsonSubCategory.Get("parentId").Exists() {
										p.Log.Debug().Msgf("No subcategory in index range(18-120): ?")
										continue
									}

									if !jsonSubCategory.Get("id").Exists() && !jsonSubCategory.Get("name").Exists() && !jsonSubCategory.Get("slag").Exists() {
										p.Log.Debug().Int("INDEX", i).Msgf("Subcategory doesn't has target property")
										continue
									}

									if jsonSubCategory.Get("parentId").String() == category.Id {
										*category.Subcategories = append(*category.Subcategories, &dto.Subcategory{
											Id:                jsonSubCategory.Get("id").String(),
											Name:              jsonSubCategory.Get("name").String(),
											Slug:              jsonSubCategory.Get("slug").String(),
											ProductCategories: &[]*dto.ProductCategory{},
										})
									}
								}
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
