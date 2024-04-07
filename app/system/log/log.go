package log

import (
	// "fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"showta.cc/app/system/conf"
	"time"
)

var lg *zap.Logger
var sugar *zap.SugaredLogger
var stdLogger *zap.SugaredLogger

func InitCore(cfg conf.Log) {
	// Custom time output format
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
	}

	// Custom log level display
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// Custom files: line number output items
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

	//Obtain encoder
	//NewJSONEncoder outputs JSON format
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = customLevelEncoder
	encoderConfig.EncodeCaller = customCallerEncoder
	//NewConsoleEncoder outputs plain text format
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	//File writeSyncer
	fullFilename := conf.AbsPath(cfg.Filename)
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fullFilename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	opts := []zapcore.WriteSyncer{}
	multiOpts := []zapcore.WriteSyncer{
		zapcore.AddSync(os.Stdout),
	}
	if cfg.Stdout {
		opts = append(opts, zapcore.AddSync(os.Stdout))
	}

	if cfg.Enable {
		opts = append(opts, fileWriteSyncer)
		multiOpts = append(multiOpts, fileWriteSyncer)
	}

	level := adapteLevel(cfg.Level)
	syncWriter := zapcore.NewMultiWriteSyncer(opts...)
	//The third and subsequent parameters are the log level for writing files, while the ErrorLevel mode only records logs at the error level
	fileCore := zapcore.NewCore(encoder, syncWriter, level)
	// log := zap.New(fileCore, zap.AddCaller())
	log := zap.New(fileCore)
	sugar = log.Sugar()

	multiSyncWriter := zapcore.NewMultiWriteSyncer(multiOpts...)
	multiCore := zapcore.NewCore(encoder, multiSyncWriter, level)
	multiLogger := zap.New(multiCore)
	stdLogger = multiLogger.Sugar()
	lg = log
}

func adapteLevel(level int) zapcore.Level {
	switch level {
	case -1:
		return zapcore.DebugLevel
	case 0:
		return zapcore.InfoLevel
	case 1:
		return zapcore.WarnLevel
	case 2:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
