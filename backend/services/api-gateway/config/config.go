package config

import "testLmd/libs/utils"

var (
	productServiceHostEnvName = "PRODUCT_SERVICE_HOST"
	productServicePortEnvName = "PRODUCT_SERVICE_PORT"
	gatewayHostEnvName        = "API_HOST"
	gatewayPortEnvName        = "API_PORT"
	cors                      = "cors"
)

type Config struct {
	HostProduct string
	PortProduct string
	HostGateway string
	PortGateway string
	Cors        string
}

func GatewayConfig() *Config {
	return &Config{
		HostProduct: utils.TrimEnv(productServiceHostEnvName),
		PortProduct: utils.TrimEnv(productServicePortEnvName),
		HostGateway: utils.TrimEnv(gatewayHostEnvName),
		PortGateway: utils.TrimEnv(gatewayPortEnvName),
		Cors:        utils.TrimEnv(cors),
	}
}
