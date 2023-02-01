package scikits

import (
	"github.com/spf13/viper"
)

var (
	MyViper = viper.New()
)

const (
	path          = "./"
	fileName      = "settings"
	fileNameLocal = "settings.local"
)

func init() {
	MyViper.AddConfigPath(path)
	MyViper.SetConfigName(fileNameLocal) //设置读取的文件名
	MyViper.SetConfigType("yaml")        //设置文件的类型

	// 判断是否有settings.local，没有的话继续 settings
	if err := MyViper.ReadInConfig(); err != nil {
		MyViper.SetConfigName(fileName) //设置读取的文件名
		if err = MyViper.ReadInConfig(); err != nil {
			panic(err)
		}
	}

}
