package config

import "fmt"

const (
	defaultConfigApiPort       = 5000
	defaultConfigApiPublicPath = "./public/"

	envConfigApiHost       = "API_HOST"
	envConfigApiPort       = "API_PORT"
	envConfigApiPublicPath = "API_PUBLIC_PATH"
)

type Api struct {
	host       string
	port       int
	publicPath string
}

func (api Api) Addr() string {
	return fmt.Sprintf("%s:%d", api.host, api.port)
}

func (api Api) PublicPath() string {
	return api.publicPath
}
