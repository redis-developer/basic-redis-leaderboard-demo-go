package config

import "fmt"

const (
	defaultConfigApiPort       = 5000
	defaultConfigApiPublicPath = "./public/"

	envConfigApiHost        = "API_HOST"
	envConfigApiPort        = "API_PORT"
	envConfigApiPublicPath  = "API_PUBLIC_PATH"
	envConfigApiTLSDisabled = "API_TLS_DISABLED"

	envExternalConfigApiPort = "PORT"
)

type Api struct {
	host        string
	port        int
	publicPath  string
	tlsDisabled bool
}

func (api Api) Addr() string {
	return fmt.Sprintf("%s:%d", api.host, api.port)
}

func (api Api) PublicPath() string {
	return api.publicPath
}

func (api Api) TLSDisabled() bool {
	return api.tlsDisabled
}
