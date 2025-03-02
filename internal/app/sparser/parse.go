package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"github.com/rs/zerolog"
	"sync"
)

func Parse(log *zerolog.Logger, geo *dto.GEO, wg *sync.WaitGroup) {
	sData := CollectSessionData(log)
	CollectProducts(sData, geo)

	wg.Done()
}
