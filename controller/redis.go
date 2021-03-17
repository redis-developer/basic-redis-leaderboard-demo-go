package controller

type Redis interface {
	HSet(key, field string, value interface{}) error
	HGetAll(key string) (map[string]string, error)
	ZAdd(key string, member string, score float64) error
	ZRevRange(key string, start, stop int64) (map[string]float64, error)
	ZRange(key string, start, stop int64) (map[string]float64, error)
	ZScore(key, member string) (float64, error)
	ZIncrBy(key string, increment float64, member string) error
	ZCount(key, min, max string) (int64, error)
}
