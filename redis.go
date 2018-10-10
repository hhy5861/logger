package logger

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"sync"
	"time"
)

var (
	entryPool = sync.Pool{
		New: func() interface{} {
			return &logrus.Entry{}
		},
	}
)

type (
	RedisConfig struct {
		Logger        *logrus.Logger
		RedisHost     string
		RedisPort     int
		RedisDB       int
		RedisKey      string
		RedisPassword string
		AppName       string
	}

	HookConfig struct {
		Key      string
		Format   string
		App      string
		Host     string
		Password string
		Hostname string
		Port     int
		DB       int
		TTL      int
	}

	RedisHook struct {
		RedisPool      *redis.Pool
		RedisHost      string
		RedisKey       string
		LogstashFormat string
		AppName        string
		Hostname       string
		RedisPort      int
		TTL            int
		Formatter      logrus.Formatter
	}

	LogstashFormatter struct {
		logrus.Formatter
		logrus.Fields
	}
)

func NewRedis(appName, redisHost, redisKey, redisPassword string, redisDB, redisPort int) *RedisConfig {
	return &RedisConfig{
		Logger:        logrus.New(),
		RedisHost:     redisHost,
		RedisKey:      redisKey,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
		RedisPort:     redisPort,
		AppName:       appName,
	}
}

func (log *RedisConfig) Output() *RedisConfig {
	hookConfig := HookConfig{
		Host:     log.RedisHost,
		Key:      log.RedisKey,
		App:      log.AppName,
		Port:     log.RedisPort,
		DB:       log.RedisDB,
		TTL:      3600,
		Password: log.RedisPassword,
	}

	log.Logger.Out = ioutil.Discard
	log.Logger.Level = logrus.InfoLevel

	hook, err := NewHook(hookConfig, DefaultFormatter(logrus.Fields{
		"type": log.AppName,
	}))

	if err == nil {
		log.Logger.AddHook(hook)
	} else {
		log.Logger.Errorf("log redis error: %q", err)
	}

	return log
}

func NewHook(config HookConfig, formatter logrus.Formatter) (*RedisHook, error) {
	pool := newRedisConnectionPool(config.Host, config.Password, config.Port, config.DB)

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to REDIS: %s", err)
	}

	return &RedisHook{
		RedisHost:      config.Host,
		RedisPool:      pool,
		RedisKey:       config.Key,
		LogstashFormat: config.Format,
		AppName:        config.App,
		Hostname:       config.Hostname,
		TTL:            config.TTL,
		Formatter:      formatter,
	}, nil

}

func (hook *RedisHook) Fire(entry *logrus.Entry) error {
	dataBytes, err := hook.Formatter.Format(entry)
	conn := hook.RedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("RPUSH", hook.RedisKey, dataBytes)
	if err != nil {
		return fmt.Errorf("error sending message to REDIS: %s", err)
	}

	if hook.TTL != 0 {
		_, err = conn.Do("EXPIRE", hook.RedisKey, hook.TTL)
		if err != nil {
			return fmt.Errorf("error setting TTL to key: %s, %s", hook.RedisKey, err)
		}
	}

	return nil
}

func (hook *RedisHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func newRedisConnectionPool(server, password string, port int, db int) *redis.Pool {
	hostPort := fmt.Sprintf("%s:%d", server, port)
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", hostPort, redis.DialDatabase(db), redis.DialPassword(password))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
