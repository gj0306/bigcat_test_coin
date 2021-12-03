package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "bigcat_test_coin/app/api"
	"bigcat_test_coin/app/blockchain/internal/config"
	"bigcat_test_coin/app/blockchain/internal/service"
	"time"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *config.WebServer, greeter *service.GreeterService,  logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{}
	//中间件
	mids := []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(logger),
	}
	if c.GrpcAddr != "" {
		opts = append(opts, grpc.Address(c.GrpcAddr))
	}
	if c.GrpcTimeout != 0 {
		opts = append(opts, grpc.Timeout(time.Second * time.Duration(c.GrpcTimeout)))
	}
	opts = append(opts, grpc.Middleware(mids...))
	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServer(srv, greeter)
	return srv
}
