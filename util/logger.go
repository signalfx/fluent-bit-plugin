package util

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(fmt string, args ...interface{})
	Info(args ...interface{})
	Infof(fmt string, args ...interface{})
	Warn(args ...interface{})
	Warnf(fmt string, args ...interface{})
	Error(args ...interface{})
	Errorf(fmt string, args ...interface{})
	Panic(args ...interface{})
	Panicf(fmt string, args ...interface{})
}

func GetLogger(logLevel zapcore.Level) Logger {
	cfg := getLoggerConfig(logLevel)

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

func getLoggerConfig(logLevel zapcore.Level) zap.Config {
	rawJSON := []byte(`{
		  "level": "` + logLevel.String() + `",
		  "encoding": "console",
		  "development": false,
		  "sampling": null,
		  "outputPaths": ["stdout"],      
		  "errorOutputPaths": ["stderr"],
	      "encoderConfig": {
	        "timeKey": "time",
	        "levelKey": "level",
	        "levelEncoder": "lowercase",
	        "messageKey": "msg"
	      }
		}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("[2006/01/02 15:04:05]"))
	}
	return cfg
}
