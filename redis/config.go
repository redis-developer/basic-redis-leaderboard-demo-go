package redis

type Config interface {
	Addr() string
	Password() string
}
