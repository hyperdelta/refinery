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
	HttpServerOptions HttpServerOptionsConfig
}

func generateRandomID(prefix string) string {
	u, _ := uuid.NewV4()
	return prefix + u.String()
}

var Debug bool = true
