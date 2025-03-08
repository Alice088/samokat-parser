package cookie

import (
	"alice088/sparser/internal/pkg/samokat"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
)

func Make(urlCookie string, cookieExpr cdp.TimeSinceEpoch) *network.SetCookieParams {
	return network.SetCookie("SELECTED_ADDRESS_KEY", urlCookie).
		WithExpires(&cookieExpr).
		WithDomain(samokat.DOMAIN).
		WithPath("/")
}
