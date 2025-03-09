package env

import (
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func GetCategoryParseLimit() int {
	limit := os.Getenv("CATEGORY_PARSE_LIMIT")

	if len(limit) == 0 {
		log.Error().Msg("CATEGORY_PARSE_LIMIT environment variable not set")
		return 2
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		log.Error().Msg("Error converting CATEGORY_PARSE_LIMIT to int")
		return 2
	}

	return limitInt
}
