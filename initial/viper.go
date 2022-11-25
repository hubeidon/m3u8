package initial

import (
	"gitee.com/don178/m3u8/global"
	"github.com/spf13/viper"
)

func init() {
	v := viper.New()
	v.SetConfigFile("conf.yaml")
	// v.AddConfigPath(".")

	err := v.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}
	err = v.Unmarshal(&global.Cfg)
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}
	global.Viper = v
}
