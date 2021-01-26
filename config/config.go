// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

//Config includes greyelk specific configurations defined in greyelk.yml
type Config struct {
	Period     time.Duration `config:"period"`
	EthRPCAddr string        `config:"eth_rpc_addr"`
	StartBlock int           `config:"start_block"`
}

//DefaultConfig sets the default values of greyelk configurations
var DefaultConfig = Config{
	Period:     1 * time.Second,
	StartBlock: -1,
}
