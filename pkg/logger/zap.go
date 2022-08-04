package logger

import (
	"go.uber.org/zap"
)

func Initialize() (logger *zap.Logger, err error) {
	logger, err = zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
