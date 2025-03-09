package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"alice088/sparser/internal/pkg/geography"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"log"
	"os"
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
			stuff := p.CollectStuff(geo)

			fileUuid := uuid.New().String()
			file, err := os.Create(fmt.Sprintf("%s.json", fileUuid))
			if err != nil {
				log.Fatalf("Ошибка при создании файла: %v", err)
			}
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
					p.Log.Error().Err(err).Str("UUID_FILE", fileUuid).Msg("Error closing output file")
				}
			}(file)

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(stuff); err != nil {
				p.Log.Error().Err(err).
					Str("REGION_ID", geo.Region).
					Str("UUID_FILE", fileUuid).
					Msg("Error encode stuff")
			}

		}()
	}

	p.wg.Wait()
}
