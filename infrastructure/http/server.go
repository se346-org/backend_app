package http

import (
	"fmt"

	"github.com/chat-socio/backend/configuration"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func NewServer(serverConfig *configuration.ServerConfig) *server.Hertz {
	// Create a new Hertz server with the specified host and port
	h := server.Default(server.WithHostPorts(fmt.Sprintf(":%d", serverConfig.Port)))
	return h
}
