package config

import (
	uuid "github.com/nu7hatch/gouuid"
)

type HttpServerOptionsConfig struct {
	ReadTimeout int
	WriteTimeout int
}

type Config struct {
	ListenAddress string
	ListenPort int
	RethinkDBAddress string
	RethinkDBPort int
	HttpServerOptions HttpServerOptionsConfig
	RethinkDBName string
}

func GenerateRandomID(prefix string) string {
	u, _ := uuid.NewV4()
	return prefix + u.String()
}

var RefineryConfig = Config {
	ListenAddress: "",
	ListenPort: 3000,
	RethinkDBAddress: "0.0.0.0",
	RethinkDBPort: 32769,
	RethinkDBName: "comsat_station",
}