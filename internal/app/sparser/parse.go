package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"github.com/rs/zerolog"
)

func Parse(log *zerolog.Logger, geo *dto.GEO) {
	sData := CollectSessionData(log)
	CollectProducts(sData)
}
