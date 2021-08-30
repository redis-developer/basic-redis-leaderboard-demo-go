package redis

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/controller"
)

const envRedisURL = "REDIS_URL"

type Value struct {
	Score float64
}

type Redis struct {
	client *redis.Client
}

func (r Redis) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(key).Result()
}

func (r Redis) HSet(key, field string, value interface{}) error {
	return r.client.HSet(key, field, value).Err()
}

func (r Redis) ZAdd(key string, member string, score float64) error {
	return r.client.ZAdd(key, redis.Z{Member: member, Score: score}).Err()
}
func (r Redis) ZRevRange(key string, start, stop int64) ([]*controller.Company, error) {
	z, err := r.client.ZRevRangeWithScores(key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	companies := make([]*controller.Company, 0, len(z))

	for i := range z {
		companies = append(companies, &controller.Company{
			Symbol:    z[i].Member.(string),
			MarketCap: z[i].Score,
		})
		//values[z[i].Member.(string)] = z[i].Score
	}

	return companies, nil

}

func (r Redis) ZRange(key string, start, stop int64) ([]*controller.Company, error) {
	z, err := r.client.ZRangeWithScores(key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	n := len(z)
	companies := make([]*controller.Company, 0, n)

	for i := range z {
		companies = append(companies, &controller.Company{
			Symbol:    z[n-i-1].Member.(string),
			MarketCap: z[n-i-1].Score,
		})
	}
	return companies, err
}

func (r Redis) ZScore(key, member string) (float64, error) {
	return r.client.ZScore(key, member).Result()
}

func (r Redis) ZIncrBy(key string, increment float64, member string) error {
	return r.client.ZIncrBy(key, increment, member).Err()
}

func (r Redis) ZCount(key, min, max string) (int64, error) {
	return r.client.ZCount(key, min, max).Result()
}

func (r Redis) Close() error {
	return r.client.Close()
}

func NewOptions(config Config) (opt *redis.Options, err error) {
	// read options from Redis URL
	url, ok := os.LookupEnv(envRedisURL)
	if ok && url != "" {
		// ref https://pkg.go.dev/github.com/go-redis/redis?utm_source=gopls#ParseURL
		opt, err = redis.ParseURL(url)
		if err != nil {
			return nil, err
		}
	} else {
		// read options from config
		opt = &redis.Options{
			Addr:     config.Addr(),
			Password: config.Password(),
		}
	}

	// read CA cert
	caPath, ok := os.LookupEnv("TLS_CA_CERT")
	if ok && caPath != "" {
		// ref https://pkg.go.dev/crypto/tls#example-Dial
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, fmt.Errorf("Failed to append root CA cert at %s", caPath)
		}
		opt.TLSConfig = &tls.Config{
			RootCAs: rootCertPool,
		}
	}

	return opt, nil
}

func New(opt *redis.Options) *Redis {
	client := redis.NewClient(opt)

	return &Redis{
		client: client,
	}
}
