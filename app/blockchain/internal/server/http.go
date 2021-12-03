package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "bigcat_test_coin/app/api"
	"bigcat_test_coin/app/blockchain/internal/service"
	"time"

	"bigcat_test_coin/app/blockchain/internal/config"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *config.WebServer, greeter *service.GreeterService,  logger log.Logger) *http.Server {
	var opts = []http.ServerOption{}
	//中间件
	mids := []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(logger),
	}
	if c.HttpAddr != "" {
		opts = append(opts, http.Address(c.HttpAddr))
	}
	if c.GrpcTimeout != 0 {
		opts = append(opts, http.Timeout(time.Second * time.Duration(c.GrpcTimeout)))
	}
	opts = append(opts, http.Middleware(mids...))
	srv := http.NewServer(opts...)
	//swagger-ui  /q/swagger-ui
	if c.Swagger{
		h := openapiv2.NewHandler()
		srv.HandlePrefix("/q/", h)
	}
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
