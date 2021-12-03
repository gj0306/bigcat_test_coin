package client

import (
	v1 "bigcat_test_coin/app/api"
	"bigcat_test_coin/app/bot/client"
)

var endpoint = "127.0.0.1:9902"
var GrpcClient *v1.GreeterClient



func NewClient() *v1.GreeterClient{
	v := client.NewOrderServiceClient(endpoint)
	GrpcClient = &v
	return GrpcClient
}

func SetClientEndpoint(addr string)  {
	endpoint = addr
	_ = NewClient()
}

func GetClient()v1.GreeterClient{
	if GrpcClient == nil{
		return *NewClient()
	}
	return *GrpcClient
}