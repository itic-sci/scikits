package scikits

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"github.com/natefinch/lumberjack"
	"log"
	"math"
	"net"
	"os"
	"path"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// MessageKey provides the key value for gelf message field
	MessageKey = "message"
	TimeKey    = "timestamp"
	CallerKey  = "caller"
)

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

func SyslogLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendInt(7)
	case zapcore.InfoLevel:
		enc.AppendInt(6)
	case zapcore.WarnLevel:
		enc.AppendInt(4)
	case zapcore.ErrorLevel:
		enc.AppendInt(3)
	case zapcore.DPanicLevel:
		enc.AppendInt(0)
	case zapcore.PanicLevel:
		enc.AppendInt(0)
	case zapcore.FatalLevel:
		enc.AppendInt(0)
	}
}

func getEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间编码器
	config.EncodeLevel = SyslogLevelEncoder
	config.MessageKey = MessageKey
	config.TimeKey = TimeKey
	config.CallerKey = CallerKey

	jsonEncode := zapcore.NewJSONEncoder(config)

	return jsonEncode
}

// GELF UDP 或 HTTP
func writeGraylogCore(host string, port int) zapcore.Core {
	allLevels := zap.LevelEnablerFunc(func(l zapcore.Level) bool { return true })
	syncer := New(NewDefaultConfig(host, port))
	jsonEncode := getEncoder()
	return zapcore.NewCore(jsonEncode, syncer, allLevels)
}

func writeFileCore() zapcore.Core {
	allLevels := zap.LevelEnablerFunc(func(l zapcore.Level) bool { return true })
	return zapcore.NewCore(getEncoder(), getLogWriter(), allLevels)
}

// Config represents the required settings for connecting the gelf data sink.
type Config struct {
	GraylogPort     int
	GraylogHostname string
	MaxChunkSize    int
}

// NewDefaultConfig provides a configuration with default values for port and chunk size.
func NewDefaultConfig(host string, port int) Config {
	return Config{GraylogPort: port, MaxChunkSize: 8154, GraylogHostname: host}
}

// New returns an implementation of ZapWriteSyncer which should be compatible with zap.WriteSyncer
func New(config Config) zapcore.WriteSyncer {
	return &gelf{Config: config}
}

type gelf struct {
	Config
}

func (g *gelf) Sync() error {
	// currently a noop.
	return nil
}

func (g *gelf) Write(p []byte) (int, error) {
	compressed, err := g.compress(p)
	if err != nil {
		return 0, err
	}
	chunksize := g.Config.MaxChunkSize
	length := compressed.Len()

	if length > chunksize {
		chunkCountInt := int(math.Ceil(float64(length) / float64(chunksize)))

		id := make([]byte, 8)
		rand.Read(id)

		for i, index := 0, 0; i < length; i, index = i+chunksize, index+1 {
			packet := g.createChunkedMessage(index, chunkCountInt, id, &compressed)
			_, e := g.send(packet.Bytes())
			if err != nil {
				return 0, e
			}
		}

	} else {
		_, e := g.send(compressed.Bytes())
		if err != nil {
			return 0, e
		}
	}

	//fmt.Printf("Wrote data: %s\n", p)
	return len(p), nil
}

func (g *gelf) createChunkedMessage(index int, chunkCountInt int, id []byte, compressed *bytes.Buffer) bytes.Buffer {
	var packet bytes.Buffer

	chunksize := g.Config.MaxChunkSize

	packet.Write(g.intToBytes(30))
	packet.Write(g.intToBytes(15))
	packet.Write(id)

	packet.Write(g.intToBytes(index))
	packet.Write(g.intToBytes(chunkCountInt))

	packet.Write(compressed.Next(chunksize))

	return packet
}

func (g *gelf) intToBytes(i int) []byte {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, int8(i))
	if err != nil {
		log.Printf("Uh oh! %s", err)
	}
	return buf.Bytes()
}

func (g *gelf) compress(b []byte) (bytes.Buffer, error) {
	// TODO enable compression
	var buf bytes.Buffer
	// comp := zlib.NewWriter(&buf)
	// defer comp.Close()
	// _, err := comp.Write(b)
	_, err := buf.Write(b)
	return buf, err
}

func (g *gelf) send(b []byte) (int, error) {
	gelfType := MyViper.GetString("graylog.gelf")
	if gelfType == "tcp" {
		return g.sendByTcp(b)
	} else {
		return g.sendByUdp(b)
	}

}

func (g *gelf) sendByUdp(b []byte) (int, error) {
	var addr = g.Config.GraylogHostname + ":" + strconv.Itoa(g.Config.GraylogPort)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Printf("Uh oh! %s", err)
		return 0, err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Printf("Uh oh! %s", err)
		return 0, err
	}
	return conn.Write(b)
}

func (g *gelf) sendByTcp(b []byte) (int, error) {
	var addr = g.Config.GraylogHostname + ":" + strconv.Itoa(g.Config.GraylogPort)
	udpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Printf("Uh oh! %s", err)
		return 0, err
	}
	conn, err := net.DialTCP("tcp", nil, udpAddr)
	if err != nil {
		log.Printf("Uh oh! %s", err)
		return 0, err
	}
	return conn.Write(b)
}

func wrapLogger(logger *zap.Logger) *zap.Logger {
	return logger.With(
		zap.Int("_pid", os.Getpid()),
		zap.String("_file", path.Base(os.Args[0])),
		zap.String("_appversion", MyViper.GetString("project.version")),
		zap.String("project", MyViper.GetString("project.name")),
	)
}
