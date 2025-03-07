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
	"time"
)

type parsingContext struct {
	log         *zerolog.Logger
	ctx         context.Context
	sessionData *dto.SessionData
}

func CollectSessionData(log *zerolog.Logger) *dto.SessionData {
	var err error
	limit := 0
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
	opts := setupChromedpOptions()

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.ListenTarget(ctx, func(event interface{}) {
			parsingContext.ctx = ctx

			switch ev := event.(type) {
			case *network.EventRequestWillBeSent:
				log.Debug().Msg(ev.Request.URL)

			case *network.EventResponseReceived:
				log.Debug().Msg(ev.Response.URL)

			}

		})

		return nil
	}))

	err = chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(samokat.MAIN),
		chromedp.Sleep(120*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return nil, &errs.ErrSessionDataMissing{}
}

func setupChromedpOptions() []chromedp.ExecAllocatorOption {
	return append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(os.Getenv("EXEC_CHROME_PATH")),
		chromedp.UserAgent(samokat.USER_AGENT),
		chromedp.Flag("headless", false),
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
