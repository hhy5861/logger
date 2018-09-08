package logger

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type KafkaHook struct {
	AppName   string
	levels    []logrus.Level
	Formatter logrus.Formatter
	producer  sarama.AsyncProducer
	Logger    *logrus.Logger
	Brokers   []string
	Topics    []string
}

func NewKafka(appName string, brokers, topics []string) *KafkaHook {
	return &KafkaHook{
		AppName: appName,
		Logger:  logrus.New(),
		Brokers: brokers,
		Topics:  topics,
	}
}

func (hook *KafkaHook) Output() *KafkaHook {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy
	kafkaConfig.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(hook.Brokers, kafkaConfig)
	if err != nil {
		return nil
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to send log entry to kafka: %v\n", err)
		}
	}()

	hook.Formatter = DefaultFormatter(logrus.Fields{
		"type": hook.AppName,
	})

	hook.Logger.WithField("topics", hook.Topics)
	hook.Logger.Hooks.Add(hook)
}

func (hook *KafkaHook) Id() string {
	return hook.AppName
}

func (hook *KafkaHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *KafkaHook) Fire(entry *logrus.Entry) error {
	var partitionKey sarama.ByteEncoder

	t, _ := entry.Data["time"].(time.Time)

	b, err := t.MarshalBinary()

	if err != nil {
		return err
	}

	partitionKey = sarama.ByteEncoder(b)

	var topics []string

	if ts, ok := entry.Data["topics"]; ok {
		if topics, ok = ts.([]string); !ok {
			return errors.New("field topics must be []string")
		}
	} else {
		return errors.New("field topics not found")
	}

	b, err = hook.Formatter.Format(entry)

	if err != nil {
		return err
	}

	value := sarama.ByteEncoder(b)
	for _, topic := range topics {
		hook.producer.Input() <- &sarama.ProducerMessage{
			Key:   partitionKey,
			Topic: topic,
			Value: value,
		}
	}

	return nil
}
