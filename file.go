package logger

import (
	"fmt"
	"os"
	"github.com/sirupsen/logrus"
	"strings"
)

type (
	FileConfig struct {
		FileName string
		SavePath string
		Logger   *logrus.Logger
		Debug    bool
	}
)

func NewFile(savePath, fileName string, debug bool) *FileConfig {
	return &FileConfig{
		Debug:    debug,
		SavePath: savePath,
		FileName: fileName,
		Logger:   logrus.New(),
	}
}

func (log *FileConfig) Output() *FileConfig {
	log.GetLoggerFullFile().CreateLogSavePath()
	log.Logger.Formatter = &logrus.JSONFormatter{}
	log.Logger.Level = logrus.InfoLevel

	file, err := os.OpenFile(log.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.Logger.Out = file
	}

	return log
}

func (log *FileConfig) CreateLogSavePath() *FileConfig {
	_, err := os.Stat(log.SavePath)
	if err != nil {
		err = os.MkdirAll(log.SavePath, os.ModePerm)
	}

	return log
}

func (log *FileConfig) GetLogPath() *FileConfig {
	if log.SavePath == "" {
		log.SavePath = "/var/log/logger"
	}

	return log
}

func (log *FileConfig) GetLogFileName() *FileConfig {
	if log.FileName == "" {
		log.FileName = "default-logger.log"
	}

	return log
}

func (log *FileConfig) GetLoggerFullFile() *FileConfig {
	log.GetLogPath().GetLogFileName()

	log.FileName = fmt.Sprintf("/%s/%s",
		strings.Trim(log.SavePath, "/"),
		strings.Trim(log.FileName, "/"))

	return log
}
