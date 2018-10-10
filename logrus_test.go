package logger

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"testing"
)

type (
	Config struct {
		Logger *Logger `yaml:"logger"`
	}
)

func TestNewLogger(t *testing.T) {
	content, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		t.Log(err)
	}

	var config *Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		t.Log(err)
	}

	// std out file
	//NewLogger(config.Logger)
	//Info("test std out file logs")

	//std out redis
	//config.Logger.StdOut = "redis"
	//NewLogger(config.Logger)
	//Info("test std out redis logs")

	//std out logstash
	//config.Logger.StdOut = "logstash"
	//NewLogger(config.Logger)
	//Info("test std out logstash logs")

	//std out kafka
	config.Logger.StdOut = "kafka"
	NewLogger(config.Logger)
	Info("test std out kafka logs")

}
