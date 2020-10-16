package util

import (
	"fmt"
	"github.com/spf13/viper"
)

//初始化配置文件
func InitConfig(){
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}