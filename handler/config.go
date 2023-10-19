package handler

type Config struct {
	RPC        string `json:"rpc" mapstructure:"rpc"`
	Resolver   string `json:"resolver" mapstructure:"resolver"`
	Registry   string `json:"registry" mapstructure:"registry"`
	PrivateKey string `json:"-" mapstructure:"privateKey"`
}