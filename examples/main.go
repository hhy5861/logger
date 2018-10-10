package main

import (
	"fmt"
	"github.com/hhy5861/logger"
	"gopkg.in/yaml.v1"
)

type (
	Config struct {
		Logger *logger.Logger `yaml:"logger"`
	}
)

var (
	content = `
logger:
  appName: test-server
  savePath: /data/logs/golang/backend-server
  stdOut: kafka    # stdout type (file|logstash|redis|kafka|elasticsearch)
  debug: true

  logStashHost: 127.0.0.1
  logStashPort: 5044

  redisHost: 127.0.0.1
  redisPort: 6379
  redisDB: 10

  brokers:
    - localhost:9092
  topics: test

  elasticHost: 127.0.0.1
  elasticPost: 9200
  prefixIndex: test
`
)

func main() {
	var config *Config
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println(config.Logger)
	logger.NewLogger(config.Logger)
	logger.Info("test std out kafka logs")

	//hook, _ := logrest.NewHook("http://localhost:9092", "test", &logrest.Options{})
	//
	//loggers := logrus.New()
	//loggers.Hooks.Add(hook)
	//
	//loggers.Info("Hello world")

}
