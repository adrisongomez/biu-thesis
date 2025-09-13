package loggers

import (
	"go.uber.org/zap"
)

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)
	return logger.Sugar(), nil
}
