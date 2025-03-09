package chrome

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

type CancelFunc = func()

func SetupContext() (context.Context, CancelFunc) {
	timeCtx, timeCtxCancel := context.WithTimeout(context.Background(), 480*time.Second)
	allocCtx, allocCancel := chromedp.NewExecAllocator(timeCtx, SetupChromedpOptions()...)
	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx)

	return chromeCtx, func() {
		defer timeCtxCancel()
		defer allocCancel()
		defer chromeCancel()
	}
}
