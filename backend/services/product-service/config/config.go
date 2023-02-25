package config

import "testLmd/libs/utils"

var (
	productServiceHostEnvName = "PRODUCT_SERVICE_HOST"
	productServicePortEnvName = "PRODUCT_SERVICE_PORT"
	pgConnStringEnvName       = "PG_CONNECTION_STRING"
)

type Config struct {
	HostProduct  string
	PortProduct  string
	PgConnString string
}

func ProductConfig() *Config {
	return &Config{
		HostProduct:  utils.TrimEnv(productServiceHostEnvName),
		PortProduct:  utils.TrimEnv(productServicePortEnvName),
		PgConnString: utils.TrimEnv(pgConnStringEnvName),
	}
}
