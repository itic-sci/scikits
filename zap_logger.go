package scikits

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
Zap提供了两种类型的日志记录器—Sugared Logger和Logger。
在性能很好但不是很关键的上下文中，使用SugaredLogger。它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。
在每一微秒和每一次内存分配都很重要的上下文中，使用Logger。它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
*/

var (
	SugarLogger *zap.SugaredLogger
	Logger      *zap.Logger
)

const (
	logLevel       = zapcore.DebugLevel
	fileMaxSize    = 50 // M
	fileMaxBackups = 10 // 文件个数
	fileMaxAge     = 30 // day
)

func init() {
	encoder := getEncoder()
	writeSyncer := getLogWriter()
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	// zap.AddCaller() 添加将调用函数信息记录到日志中的功能。
	Logger = zap.New(core, zap.AddCaller())
	SugarLogger = Logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间编码器

	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	//return zapcore.NewConsoleEncoder(encoderConfig) // printF输出：2022-11-07T10:21:06.970+0800	ERROR	test_main/test_zap_log.go:41	runtime error: integer divide by zero

	return zapcore.NewJSONEncoder(encoderConfig) // json 格式输出：{"level":"ERROR","ts":"2022-11-07T10:26:59.189+0800","caller":"test_main/test_zap_log.go:41","msg":"runtime error: integer divide by zero"}
}

func getLogWriter() zapcore.WriteSyncer {
	logFile := MyViper.GetString("logs.filepath")
	if logFile == "" {
		logFile = "./logs/zap_logger.log"
	}

	// 日志文件每50MB会切割并且在当前目录下最多保存5个备份
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    fileMaxSize,
		MaxBackups: fileMaxBackups,
		MaxAge:     fileMaxAge,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
