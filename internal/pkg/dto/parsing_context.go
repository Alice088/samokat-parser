package dto

import "context"

type ParsingContext struct {
	Ctx  context.Context
	Skip *map[string]bool
}
