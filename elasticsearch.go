package logger

import (
	"fmt"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
	"time"
)

type (
	ElasticConfig struct {
		ElasticHost string
		ElasticPost int
		PrefixIndex string
		Logger      *logrus.Logger
	}
)

func NewElastic(host, prefixIndex string, post int) *ElasticConfig {
	return &ElasticConfig{
		ElasticHost: host,
		ElasticPost: post,
		PrefixIndex: prefixIndex,
		Logger:      logrus.New(),
	}
}

func (esc *ElasticConfig) Output() *ElasticConfig {
	address := fmt.Sprintf("http://%s:%d", esc.ElasticHost, esc.ElasticPost)
	client, err := elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		logrus.Fatal(err)
	}

	hook, err := elogrus.NewAsyncElasticHookWithFunc(client, esc.ElasticHost, logrus.DebugLevel, func() string {
		t := time.Now()
		index := fmt.Sprintf("%s-%d-%d-%d", esc.PrefixIndex, t.Year(), t.Month(), t.Day())
		return index
	})

	if err != nil {
		logrus.Fatal(err)
	}

	esc.Logger.Hooks.Add(hook)

	return esc
}
