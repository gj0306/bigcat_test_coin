package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kgtpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"google.golang.org/grpc"
	v1 "bigcat_test_coin/app/api/contapi"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService,NewServiceClient)

// NewServiceClient 服务发现Client
func NewServiceClient(endpoint string) v1.SourceGreeterClient {
	if endpoint == ""{
		endpoint = "127.0.0.1:9902"
	}
	conn := NewGrpcClientConn(endpoint)
	return v1.NewSourceGreeterClient(conn)
}

func NewGrpcClientConn(endpoint string)*grpc.ClientConn{
	mids := [] middleware.Middleware{
		recovery.Recovery(),
	}
	conn, err := kgtpc.DialInsecure(
		context.Background(),
		kgtpc.WithEndpoint(endpoint),
		kgtpc.WithMiddleware(mids...),
	)
	if err != nil {
		panic(err)
	}
	return conn
}