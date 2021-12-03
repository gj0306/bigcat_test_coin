package client

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kgtpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	v1 "bigcat_test_coin/app/api"
)

// NewOrderServiceClient 服务发现Client
func NewOrderServiceClient(endpoint string) v1.GreeterClient {
	if endpoint == ""{
		endpoint = "127.0.0.1:9902"
	}
	conn := NewGrpcClientConn(endpoint)
	return v1.NewGreeterClient(conn)
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