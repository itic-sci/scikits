package main

import (
	"github.com/itic-sci/scikits"
)

func main() {
	defer scikits.SugarLogger.Sync()
	scikits.SugarLogger.Info("udp info 测试")
	//scikits.SugarLogger.Error("1231231231312 徐威 error message")
}
