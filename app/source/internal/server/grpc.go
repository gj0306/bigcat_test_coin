package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "bigcat_test_coin/app/api/contapi"
	"bigcat_test_coin/app/source/internal/service"
	"time"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(greeter *service.GreeterService, logger log.Logger,) *grpc.Server {
	var opts = []grpc.ServerOption{}
	//中间件
	mids := []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(logger),
	}

	opts = append(opts, grpc.Network(blockChainGrpc))
	opts = append(opts, grpc.Timeout(time.Second*2))

	opts = append(opts, grpc.Middleware(mids...))
	srv := grpc.NewServer(opts...)
	v1.RegisterSourceGreeterServer(srv, greeter)
	return srv
}
