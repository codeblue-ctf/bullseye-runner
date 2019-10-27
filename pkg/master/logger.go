package master

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger(debug bool) {
	var err error
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("failed to init logger: %+v", err)
	}
}
