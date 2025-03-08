package dto

import (
	"alice088/sparser/internal/pkg/chrome"
	"alice088/sparser/internal/pkg/cookie"
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"sync"
	"time"
)

type ParsingContext struct {
	ChromeCtx context.Context
	EventCtx  context.Context
	Skip      *sync.Map
	GeoCookie *network.SetCookieParams
}

func NewParsingContext(skip *sync.Map, urlCookie string) (*ParsingContext, chrome.CancelFunc) {
	chromeCtx, chromeCancel := chrome.SetupContext()

	return &ParsingContext{
		Skip:      skip,
		GeoCookie: cookie.Make(urlCookie, cdp.TimeSinceEpoch(time.Now().Add(180*24*time.Hour))),
		ChromeCtx: chromeCtx,
	}, chromeCancel
}
