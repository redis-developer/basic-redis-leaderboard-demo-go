package config

import "fmt"

const (
	defaultConfigRedisPort = 6379

	envConfigRedisHost     = "REDIS_HOST"
	envConfigRedisPort     = "REDIS_PORT"
	envConfigRedisPassword = "REDIS_PASSWORD"
)

type Redis struct {
	host     string
	port     int
	password string
}

func (redis Redis) Addr() string {
	return fmt.Sprintf("%s:%d", redis.host, redis.port)
}

func (redis Redis) Password() string {
	return redis.password
}
