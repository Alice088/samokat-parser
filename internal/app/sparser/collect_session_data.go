package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"alice088/sparser/internal/pkg/samokat"
	"context"
	"encoding/json"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

type parsingContext struct {
	ev          *network.EventRequestWillBeSent
	ctx         context.Context
	sessionData *dto.SessionData
	log         *zerolog.Logger
}

func CollectSessionData(log *zerolog.Logger) *dto.SessionData {
	sessionData := new(dto.SessionData)
	parsCtx := &parsingContext{
		log:         log,
		sessionData: sessionData,
	}
	var skipEvent bool

	opts := setupChromedpOptions()

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.ListenTarget(ctx, func(ev interface{}) {
			if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
				parsCtx.ev = ev
				parsCtx.ctx = ctx

				if strings.Contains(ev.Request.URL, samokat.CARTS_GET) && ev.Request.HasPostData && !skipEvent {
					log.Debug().Msg(ev.Request.URL)
					collectAuthToken(parsCtx)
					go collectShowcaseId(parsCtx)
					skipEvent = true
				}
			}
		})

		return nil
	}))

	err = chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(samokat.MAIN),
		chromedp.Sleep(5*time.Second),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Fatal during chrome run")
	}

	return sessionData
}

func setupChromedpOptions() []chromedp.ExecAllocatorOption {
	return append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(os.Getenv("EXEC_CHROME_PATH")),
		chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"),
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)
}

func collectShowcaseId(parsCtx *parsingContext) {
	postData, err := network.GetRequestPostData(parsCtx.ev.RequestID).Do(parsCtx.ctx)

	err = json.Unmarshal([]byte(postData), parsCtx.sessionData)
	if err != nil {
		parsCtx.log.Error().Err(err).Msg("Failed to get showcaseId")
	}

	parsCtx.log.Debug().Str("ShowcaseId", parsCtx.sessionData.ShowcaseID).Send()
}

func collectAuthToken(parsCtx *parsingContext) {
	if authTkn, ok := parsCtx.ev.Request.Headers["authorization"].(string); ok {
		parsCtx.sessionData.AuthToken = authTkn
	}
	parsCtx.log.Debug().Any("Token", parsCtx.sessionData.AuthToken).Send()
}
