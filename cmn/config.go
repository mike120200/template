package cmn

import (
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// ViperInit 初始化viper
func ViperInit(configName string) error {
	// 1. 加载环境变量
	if err := gotenv.Load(); err != nil {
		log.Println(".env file not found, skip loading it")
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 优先使用环境变量
	viper.AutomaticEnv()

	viper.SetConfigName(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.AddConfigPath("../../..")
	viper.AddConfigPath("../../../..")
	viper.AddConfigPath("../../../../..")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
