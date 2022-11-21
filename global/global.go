package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)


var (
	Log *zap.Logger
	Slog *zap.SugaredLogger
	Viper *viper.Viper
)