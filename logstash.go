package logger

import (
	"github.com/sirupsen/logrus"
	"net"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"fmt"
)

type (
	LogStashConfig struct {
		Logger       *logrus.Logger
		LogStashHost string
		LogStaShPort int
		AppName      string
	}
)

func NewLogStash(appName, logStashHost string, logStaShPort int) *LogStashConfig {
	return &LogStashConfig{
		Logger:       logrus.New(),
		AppName:      appName,
		LogStashHost: logStashHost,
		LogStaShPort: logStaShPort,
	}
}

func (log *LogStashConfig) Output() *LogStashConfig {
	conn, err := net.Dial("tcp", log.getLogStashAddress())
	if err != nil {
		log.Logger.Fatal(err)
	}

	log.Logger.Formatter = &logrus.JSONFormatter{}
	log.Logger.Level = logrus.InfoLevel

	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": log.AppName}))

	log.Logger.Hooks.Add(hook)

	return log
}

func (log *LogStashConfig) getLogStashAddress() string {
	return fmt.Sprintf("%s:%d", log.LogStashHost, log.LogStaShPort)
}
