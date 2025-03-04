package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

func InitLogger() *zap.Logger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "D:/goland/Gavna/logs/app.log", //сомнительно
		MaxSize:    50,
		MaxBackups: 5,
		MaxAge:     28,
		Compress:   true,
	}

	fileWriter := zapcore.AddSync(lumberjackLogger)
	consoleWriter := zapcore.AddSync(os.Stdout)
	multiWriter := zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		multiWriter,
		zap.InfoLevel,
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return Logger
}
