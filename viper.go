package scikits

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

var (
	MyViper       = viper.New()
	MyViperConfig = ""
)

const (
	defaultFile  = "settings.local.yaml"
	prodFileName = "settings.yaml"
)

func init() {
	flag.StringVar(&MyViperConfig, "c", "", "choose config file.")
	flag.Parse()

	if MyViperConfig == "" { // 判断命令行参数是否为空
		MyViperConfig = defaultFile
		fmt.Printf("您正在使用默认的config路径，为%s\n", MyViperConfig)
	} else { // 命令行参数不为空 将值赋值于config
		fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", MyViperConfig)
	}

	MyViper.SetConfigFile(MyViperConfig) //设置读取的文件名
	MyViper.SetConfigType("yaml")        //设置文件的类型

	// 判断是否有settings.local，没有的话继续 settings
	if err := MyViper.ReadInConfig(); err != nil {
		MyViperConfig = prodFileName
		fmt.Printf("您正在使用默认的config路径，为%s\n", MyViperConfig)
		MyViper.SetConfigFile(MyViperConfig) //设置读取文件
		if err = MyViper.ReadInConfig(); err != nil {
			panic(err)
		}
	}
}
