package log

import "go.uber.org/zap"

var Logger *zap.Logger

func Init(prod bool) error {
	var err error
	if prod {
		Logger, err = zap.NewProduction()
	} else {
		Logger, err = zap.NewDevelopment()
	}
	return err
}
