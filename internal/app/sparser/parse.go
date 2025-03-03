package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"github.com/rs/zerolog"
	"sync"
)

func Parse(log *zerolog.Logger, geo *dto.GEO, wg *sync.WaitGroup) {
	sessionData := CollectSessionData(log)
	categories := CollectCategories(log, sessionData, geo)

	log.Debug().Interface("categories", categories).Msg("parsed categories")
	wg.Done()
}
