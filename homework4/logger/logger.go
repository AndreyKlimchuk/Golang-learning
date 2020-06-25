package logger

import "go.uber.org/zap"

var Zap *zap.Logger

func InitZap() (err error) {
	Zap, err = zap.NewProduction()
	return err
}
