package api

type Config interface {
	Addr() string
	PublicPath() string
}
