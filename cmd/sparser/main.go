package main

import (
	"alice088/sparser/internal/app/sparser"
	"alice088/sparser/internal/pkg/env"
	"alice088/sparser/internal/pkg/logger"
)

func main() {
	env.Init()
	log, logFile := logger.Init()
	defer logger.CloseLog(logFile)
	parser := sparser.NewParser(log)

	parser.Parse()
}
