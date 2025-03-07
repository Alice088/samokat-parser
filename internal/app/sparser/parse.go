package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"alice088/sparser/internal/pkg/geography"
	"github.com/rs/zerolog"
	"sync"
)

type Parser struct {
	Log  *zerolog.Logger
	Geos []*dto.GEO
	wg   *sync.WaitGroup
}

func NewParser(log *zerolog.Logger) *Parser {
	geos, err := geography.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing geos")
	}

	group := sync.WaitGroup{}
	group.Add(len(geos))

	return &Parser{
		Log:  log,
		wg:   &group,
		Geos: geos,
	}
}

func (p *Parser) Parse() {
	for _, geo := range p.Geos {
		p.Log.Debug().Int("GEO_ID", geo.ID).Msgf("Current parsing region: %s", geo.Region)

		go func() {
			defer p.wg.Done()
			_, err := p.CollectStuff(geo)
			if err != nil {
				p.Log.Error().Err(err).Str("Region", geo.Region).Msg("Error collecting stuff")
			}
		}()
	}

	p.wg.Wait()
}
