package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/braintree/manners"
)

type HttpServerConfig struct {
	Port     uint          `mapstructure:"port"`
	Hostname string        `mapstructure:"hostname"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type HttpServerAdapter struct {
	server *manners.GracefulServer
	config HttpServerConfig
}

func (hsa *HttpServerAdapter) Close() error {

	if hsa.server.BlockingClose() {
		return nil
	}

	return fmt.Errorf("action Close already called")
}

func (hsa *HttpServerAdapter) ListenAndServe() error {

	return hsa.server.ListenAndServe()
}

func (hsa *HttpServerAdapter) SetHandler(handler http.Handler) {

	hsa.server.Handler = handler
}

func NewHttpServerAdapter(config HttpServerConfig) HttpServerAdapter {

	httpServer := http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Hostname, config.Port),
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
	}

	server := manners.NewWithServer(&httpServer)

	return HttpServerAdapter{
		config: config,
		server: server,
	}
}
