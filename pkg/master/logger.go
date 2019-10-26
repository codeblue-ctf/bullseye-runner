package master

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger: %+v", err)
	}
}
