package main

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

func generateRandomNodeID() string {
	u, _ := uuid.NewV4()
	return "solo-" + u.String()
}
