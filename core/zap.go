package core

import (
	"log"
	"os"
	"server/global"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.Logger {
	zapConfig := global.Config.Zap

	writerSyncer := getLoggerWriter(zapConfig.Filename, zapConfig.MaxSize, zapConfig.MaxBackups, zapConfig.MaxAge)
	if zapConfig.IsConsolePrint {
		writerSyncer = zapcore.NewMultiWriteSyncer(writerSyncer, zapcore.AddSync(os.Stdout))
	}
	encoder := getEncoder()

	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(zapConfig.Level))
	if err != nil {
		log.Fatalf("Failed to parse log level:%v", err)
	}
	core := zapcore.NewCore(encoder, writerSyncer, logLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

func getLoggerWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
