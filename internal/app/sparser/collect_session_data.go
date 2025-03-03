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
	var sData *dto.SessionData
	var err error
	limit := 2
	count := 0

	for {
		sData, err = runCollect(log)

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

		return sData
	}
}

func runCollect(log *zerolog.Logger) (*dto.SessionData, error) {
	var skipEvent bool
	opts := setupChromedpOptions()
	sessionData := new(dto.SessionData)

	parsCtx := &parsingContext{
		log:         log,
		sessionData: sessionData,
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.ListenTarget(ctx, func(ev interface{}) {
			if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
				parsCtx.ev = ev
				parsCtx.ctx = ctx

				if strings.Contains(ev.Request.URL, samokat.CARTS_GET) && ev.Request.HasPostData && !skipEvent {
					log.Debug().Msg(ev.Request.URL)
					log.Debug().Msg(string(ev.RequestID))
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
		return nil, err
	}

	if sessionData.AuthToken != "" && sessionData.ShowcaseID != "" {
		return sessionData, nil
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
