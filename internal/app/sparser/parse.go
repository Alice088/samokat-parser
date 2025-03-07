package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"alice088/sparser/internal/pkg/geography"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type Parser struct {
	Log *zerolog.Logger
	Geo []*dto.GEO
	wg  *sync.WaitGroup
}

func NewParser(log *zerolog.Logger) *Parser {
	geo, err := geography.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing geo")
	}

	group := sync.WaitGroup{}
	group.Add(len(geo))

	return &Parser{
		Log: log,
		wg:  &group,
		Geo: geo,
	}
}

func (p *Parser) Parse() {
	for _, geo := range p.Geo {
		log.Debug().Int("GEO_ID", geo.ID).Msgf("Current parsing region: %s", geo.Region)

		go func() {
			defer p.wg.Done()
			CollectSessionData()
		}()
	}

	p.wg.Wait()
}
