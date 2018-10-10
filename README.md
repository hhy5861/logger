#### logrus std out (file | kafka | logstash | redis | elasticsearch)

##
```yaml
std out file config:

logger:
    stdOut: file
    debug: true
    savePath: /data/logs/golang/backend-server
```

##
```yaml
std out logstash config:

logger:
    stdOut: logstash
    debug: true
    logStashHost: 127.0.0.1
    logStashPort: 30693
```

##
```yaml
std out redis config:

logger:
    stdOut: redis
    debug: true
    redisHost: 127.0.0.1
    redisPort: 30079
    redisDB: 10
```

##
```yaml
std out kafka config:

logger:
    stdOut: kafka
    debug: true
    brokers: 
      - 127.0.0.1
      - 10.10.100.21
    topics: 
      - errors
      - info
      - warn
```

use:
```go

```