package global

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Address struct {
	Path   string `json:"path,omitempty" yaml:"path"`
	Prefix string `json:"prefix,omitempty" yaml:"prefix"`
	Fname string `json:"fname,omitempty" yaml:"fname"`
}

type Config struct {
	Dir            string        `json:"dir,omitempty" yaml:"dir"`
	UserAgent      string        `json:"user_agent,omitempty" yaml:"userAgent"`
	RequestTimeout time.Duration `json:"request_timeout,omitempty" yaml:"requestTimeout"`
	Address        []Address     `json:"address,omitempty" yaml:"address"`
	Ext            string        `json:"ext,omitempty" yaml:"ext"`
}

var (
	Log   *zap.Logger
	Slog  *zap.SugaredLogger
	Viper *viper.Viper
	Cfg   Config
)
