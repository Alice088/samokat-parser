package main

import (
	"alice088/sparser/internal/app/sparser"
	"alice088/sparser/internal/pkg/env"
	"alice088/sparser/internal/pkg/geo"
	"alice088/sparser/internal/pkg/logger"
)

func main() {
	env.Init()
	log, logFile := logger.Init()
	geos, err := geo.Init()
	defer logger.CloseLog(logFile)

	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing geo")
	}

	for _, geoDto := range geos {
		log.Debug().Interface("geo", geoDto).Msg("Geo")
		sparser.Parse(log, geoDto)
	}

}
