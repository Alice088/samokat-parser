package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	errs "alice088/sparser/internal/pkg/errors"
	"errors"
	"github.com/rs/zerolog"
	"sync"
)

func Parse(log *zerolog.Logger, geo *dto.GEO, wg *sync.WaitGroup) {
	var sData *dto.SessionData
	var err error
	limit := 2
	count := 0

	for {
		sData, err = CollectSessionData(log)

		if err != nil {
			if errors.Is(err, &errs.ErrSessionDataMissing{}) && count < limit {

				log.Error().Err(err).
					Int("Try", count).
					Msgf("Error session data missing for: %s. Starting new try", geo.Region)

				count++
				continue
			}

			log.Error().Err(err).Msg("Error during CollectSessionData")
			break
		}

		break
	}

	CollectProducts(sData, geo)
	wg.Done()
}
