package logger

import (
	"runtime"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type (
	Logger struct {
		Debug         bool     `yaml:"debug"`
		StdOut        string   `yaml:"stdOut"`
		AppName       string   `yaml:"appName"`
		FileName      string   `yaml:"fileName"`
		SavePath      string   `yaml:"savePath"`
		LogStashHost  string   `yaml:"logStashHost"`
		LogStashPort  int      `yaml:"logStashPort"`
		RedisHost     string   `yaml:"redisHost"`
		RedisPort     int      `yaml:"redisPort"`
		RedisDB       int      `yaml:"redisDB"`
		RedisKey      string   `yaml:"redisKey"`
		RedisPassword string   `yaml:"redisPassword"`
		Brokers       []string `yaml:"brokers"`
		Topics        []string `yaml:"topics"`
	}
)

var (
	LLogger *logrus.Logger
)

func NewLogger(logger *Logger) {
	stdOut := strings.ToLower(logger.StdOut)
	switch stdOut {
	case "logstash":
		logger := NewLogStash(logger.AppName, logger.LogStashHost, logger.LogStashPort).Output()
		LLogger = logger.Logger
		break
	case "redis":
		logger := NewRedis(logger.AppName, logger.RedisHost, logger.RedisKey, logger.RedisPassword, logger.RedisDB, logger.RedisPort).Output()
		LLogger = logger.Logger
		break
	case "kafka":
		logger := NewKafka(logger.AppName, logger.Brokers, logger.Topics)
		LLogger = logger.Logger
	default:
		logger := NewFile(logger.SavePath, logger.FileName, logger.Debug).Output()
		LLogger = logger.Logger
		break
	}
}

func GetLogger() *logrus.Logger {
	return LLogger
}

func Info(message ... interface{}) {
	_, file, line, _ := runtime.Caller(1)
	files := fmt.Sprintf("%s (%d)", file, line)

	LLogger.WithFields(logrus.Fields{
		"files": files,
	}).Info(message)
}

func Warn(err error, message ... interface{}) {
	_, file, line, _ := runtime.Caller(1)
	files := fmt.Sprintf("%s (%d)", file, line)

	LLogger.WithFields(logrus.Fields{
		"files":  files,
		"errors": err,
	}).Warn(message)
}

func Fatal(err error, message ... interface{}) {
	_, file, line, _ := runtime.Caller(1)
	files := fmt.Sprintf("%s (%d)", file, line)

	LLogger.WithFields(logrus.Fields{
		"files":  files,
		"errors": err,
	}).Fatal(message)
}

func Error(err error, message ... interface{}) {
	_, file, line, _ := runtime.Caller(1)
	files := fmt.Sprintf("%s (%d)", file, line)

	LLogger.WithFields(logrus.Fields{
		"files":  files,
		"errors": err,
	}).Error(message)

}

func Debug(message ... interface{}) {
	_, file, line, _ := runtime.Caller(1)
	files := fmt.Sprintf("%s (%d)", file, line)

	LLogger.WithFields(logrus.Fields{
		"files": files,
	}).Debug(message)
}
