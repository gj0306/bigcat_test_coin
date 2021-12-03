package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	v1 "bigcat_test_coin/app/api/contapi"
	"bigcat_test_coin/app/source/internal/service"
	"time"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{}
	//中间件
	mids := []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(logger),
	}
	opts = append(opts, http.Network(blockChainHttp))
	opts = append(opts, http.Timeout(time.Second*2))

	opts = append(opts, http.Middleware(mids...))
	srv := http.NewServer(opts...)
	//swagger-ui  /q/swagger-ui
	h := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", h)

	v1.RegisterSourceGreeterHTTPServer(srv, greeter)
	return srv
}
