package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	errs "alice088/sparser/internal/pkg/errors"
	"alice088/sparser/internal/pkg/samokat"
	"context"
	"encoding/json"
	"errors"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

type parsingContext struct {
	log         *zerolog.Logger
	ctx         context.Context
	sessionData *dto.SessionData
}

func CollectSessionData(log *zerolog.Logger) *dto.SessionData {
	var err error
	limit := 2
	count := 0

	parsContext := &parsingContext{
		log: log,
		sessionData: &dto.SessionData{
			AuthToken:  "",
			ShowcaseID: "",
		},
	}

	for {
		parsContext.sessionData, err = runCollect(log, parsContext)

		if err != nil {
			if errors.Is(err, &errs.ErrSessionDataMissing{}) && count < limit {
				log.Error().Err(err).
					Int("Try", count).
					Msg("Error session data missing for. Starting new try")

				count++
				continue
			}

			log.Fatal().Err(err).Msg("Error during CollectSessionData")
		}

		return parsContext.sessionData
	}
}

func runCollect(log *zerolog.Logger, parsingContext *parsingContext) (*dto.SessionData, error) {
	skip := map[string]bool{
		samokat.CARTS_GET:       false,
		samokat.CATEGORIES_LIST: false,
	}

	rand.Seed(time.Now().UnixNano())

	opts := setupChromedpOptions()

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.ListenTarget(ctx, func(event interface{}) {
			parsingContext.ctx = ctx

			switch ev := event.(type) {
			case *network.EventRequestWillBeSent:
				if strings.Contains(ev.Request.URL, samokat.CARTS_GET) && ev.Request.HasPostData && !skip[samokat.CARTS_GET] {
					log.Debug().Msg(ev.Request.URL)
					log.Debug().Msg(string(ev.RequestID))

					collectAuthToken(ev, parsingContext)
					go collectShowcaseId(ev, parsingContext)
					skip[samokat.CARTS_GET] = true
				}

			case *network.EventResponseReceived:
				log.Debug().Msg(ev.Response.URL)

				if strings.Contains(ev.Response.URL, "https://api-web.samokat.ru/v2/showcases/35ad943e-c956-4a31-9333-f5e02cc5c2b0/categories/0549e93d-444a-4975-8861-8f15a2285724") {
					time.Sleep(3 * time.Second)
					log.Debug().Msg(ev.Response.URL)

					go func() {
						body, err := network.GetResponseBody(ev.RequestID).Do(parsingContext.ctx)
						if err != nil {
							log.Error().Err(err).Msg("Error getting response body")
						}
						log.Debug().Msg(string(body))

					}()
				}

				if match := regexp.MustCompile(samokat.CATEGORIES_LIST_REGEXP).FindStringSubmatch(ev.Response.URL); match != nil {
					log.Debug().Msg(ev.Response.URL)
					log.Debug().Msg(string(ev.RequestID))

					time.Sleep(time.Second)
					go func() {
						body, err := network.GetResponseBody(ev.RequestID).Do(parsingContext.ctx)

						if len(body) == 0 && err != nil {
							return
						}

						if err != nil {
							log.Error().Err(err).Msg("Error getting response body")
						}

						result := gjson.Get(string(body), "0")

						log.Debug().Msg(string(result.String()))

					}()
				}

			}

		})

		return nil
	}))

	var screen []byte
	headers := map[string]string{
		"Authorization": "Bearer " + parsingContext.sessionData.AuthToken, // Добавление токена в заголовок
		"User-Agent":    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	}

	err = chromedp.Run(ctx,
		network.Enable(),
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(samokat.MAIN),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for key, value := range headers {
				if err = network.SetExtraHTTPHeaders(map[string]any{key: value}).Do(ctx); err != nil {
					return err
				}
			}
			return nil
		}),
		chromedp.Sleep(10*time.Second),
		chromedp.Navigate("https://samokat.ru/category/torty-i-pirozhnye-2"),
		chromedp.Sleep(4*time.Second),
		chromedp.CaptureScreenshot(&screen),
		chromedp.Sleep(10*time.Second),
	)

	err = os.WriteFile("screenshot.png", screen, 0644)
	if err != nil {
		//log.Fatal(err)
	}
	if err != nil {
		return nil, err
	}

	if parsingContext.sessionData.AuthToken != "" && parsingContext.sessionData.ShowcaseID != "" {
		return parsingContext.sessionData, nil
	}

	return nil, &errs.ErrSessionDataMissing{}
}

func setupChromedpOptions() []chromedp.ExecAllocatorOption {
	return append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(os.Getenv("EXEC_CHROME_PATH")),
		chromedp.UserAgent(samokat.USER_AGENT),
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)
}

func collectShowcaseId(ev *network.EventRequestWillBeSent, parsingContext *parsingContext) {
	postData, err := network.GetRequestPostData(ev.RequestID).Do(parsingContext.ctx)

	if err != nil {
		parsingContext.log.Error().Err(err).Msg("Failed to get request post data")
	}

	if err = json.Unmarshal([]byte(postData), parsingContext.sessionData); err != nil {
		parsingContext.log.Error().Err(err).Msg("Failed to get showcaseId")
	} else {
		parsingContext.log.Debug().Str("ShowcaseId", parsingContext.sessionData.ShowcaseID).Send()
	}
}

func collectAuthToken(ev *network.EventRequestWillBeSent, parsingContext *parsingContext) {
	if authTkn, ok := ev.Request.Headers["authorization"].(string); ok {
		parsingContext.sessionData.AuthToken = authTkn
		parsingContext.log.Debug().Any("Token", authTkn).Send()
	}
}
