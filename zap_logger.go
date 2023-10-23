package scikits

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
Zap提供了两种类型的日志记录器—Sugared Logger和Logger。
在性能很好但不是很关键的上下文中，使用SugaredLogger。它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。
在每一微秒和每一次内存分配都很重要的上下文中，使用Logger。它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
*/

const (
	logLevel       = zapcore.DebugLevel
	fileMaxSize    = 50 // M
	fileMaxBackups = 10 // 文件个数
	fileMaxAge     = 30 // day
	appVersion     = "0.0.1"
)

var (
	SugarLogger *zap.SugaredLogger
	Logger      *zap.Logger
)

func init() {
	var core zapcore.Core
	if MyViper.GetString("graylog.host") != "" {
		core = writeGraylogCore(MyViper.GetString("graylog.host"), MyViper.GetInt("graylog.port"))
	} else {
		core = writeFileCore()
	}

	// zap.AddCaller() 添加将调用函数信息记录到日志中的功能。
	Logger = wrapLogger(zap.New(core, zap.AddCaller()))
	SugarLogger = Logger.Sugar()
}
