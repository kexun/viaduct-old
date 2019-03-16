package config

import "flag"

type Config struct {
	Addr string
	CmdType string
	KeyFile string
	CertFile string
	CaFile string
}

func InitConfig() *Config {
	config := Config{}

	flag.StringVar(&config.Addr, "addr", "", "client or server ip addr")
	flag.StringVar(&config.CmdType, "type", "", "client or server")
	flag.StringVar(&config.KeyFile, "key", "", "key file path")
	flag.StringVar(&config.CertFile, "cert", "", "cert file path")
	flag.StringVar(&config.CaFile, "ca", "", "ca file path")

	flag.Parse()
	return &config
}