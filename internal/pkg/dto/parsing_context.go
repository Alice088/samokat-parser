package dto

import "context"

type ParsingContext struct {
	ChromeCtx context.Context
	EventCtx  context.Context
	Skip      *map[string]bool
}
