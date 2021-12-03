// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.0.3

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type GreeterHTTPServer interface {
	SayAccount(context.Context, *GetAccountRequest) (*GetAccountReply, error)
	SayBlock(context.Context, *BlockRequest) (*BlockReply, error)
	SayBlocks(context.Context, *BlocksRequest) (*BlocksReply, error)
	SayCont(context.Context, *GetContRequest) (*GetContReply, error)
	SayCreateCont(context.Context, *CreateContRequest) (*CreateContReply, error)
	SayCreateTransaction(context.Context, *CreateTransactionRequest) (*CreateTransactionReply, error)
	SayMiners(context.Context, *GetMinersRequest) (*GetMinersReply, error)
	SayNodes(context.Context, *GetNodesRequest) (*GetNodesReply, error)
	SayTransaction(context.Context, *GetTransactionRequest) (*GetTransactionReply, error)
}

func RegisterGreeterHTTPServer(s *http.Server, srv GreeterHTTPServer) {
	r := s.Route("/")
	r.GET("/blocks", _Greeter_SayBlocks0_HTTP_Handler(srv))
	r.GET("/block/{parm}", _Greeter_SayBlock0_HTTP_Handler(srv))
	r.GET("/transaction/{tx}", _Greeter_SayTransaction0_HTTP_Handler(srv))
	r.GET("/cont/{tx}", _Greeter_SayCont0_HTTP_Handler(srv))
	r.GET("/account/{addr}", _Greeter_SayAccount0_HTTP_Handler(srv))
	r.GET("/miners", _Greeter_SayMiners0_HTTP_Handler(srv))
	r.GET("/nodes", _Greeter_SayNodes0_HTTP_Handler(srv))
	r.POST("/transaction", _Greeter_SayCreateTransaction0_HTTP_Handler(srv))
	r.POST("/cont", _Greeter_SayCreateCont0_HTTP_Handler(srv))
}

func _Greeter_SayBlocks0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in BlocksRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/pro.v1.Greeter/SayBlocks")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayBlocks(ctx, req.(*BlocksRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*BlocksReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_SayBlock0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in BlockRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/pro.v1.Greeter/SayBlock")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayBlock(ctx, req.(*BlockRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*BlockReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_SayTransaction0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetTransactionRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/pro.v1.Greeter/SayTransaction")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayTransaction(ctx, req.(*GetTransactionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetTransactionReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_SayCont0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetContRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/pro.v1.Greeter/SayCont")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayCont(ctx, req.(*GetContRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetContReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_SayAccount0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetAccountRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/pro.v1.Greeter/SayAccount")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayAccount(ctx, req.(*GetAccountRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetAccountReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_SayMiners0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetMinersRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/pro.v1.Greeter/SayMiners")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayMiners(ctx, req.(*GetMinersRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetMinersReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_SayNodes0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetNodesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/pro.v1.Greeter/SayNodes")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayNodes(ctx, req.(*GetNodesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetNodesReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_SayCreateTransaction0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateTransactionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/pro.v1.Greeter/SayCreateTransaction")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayCreateTransaction(ctx, req.(*CreateTransactionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateTransactionReply)
		return ctx.Result(200, reply)
	}
}

func _Greeter_SayCreateCont0_HTTP_Handler(srv GreeterHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateContRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/pro.v1.Greeter/SayCreateCont")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SayCreateCont(ctx, req.(*CreateContRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateContReply)
		return ctx.Result(200, reply)
	}
}

type GreeterHTTPClient interface {
	SayAccount(ctx context.Context, req *GetAccountRequest, opts ...http.CallOption) (rsp *GetAccountReply, err error)
	SayBlock(ctx context.Context, req *BlockRequest, opts ...http.CallOption) (rsp *BlockReply, err error)
	SayBlocks(ctx context.Context, req *BlocksRequest, opts ...http.CallOption) (rsp *BlocksReply, err error)
	SayCont(ctx context.Context, req *GetContRequest, opts ...http.CallOption) (rsp *GetContReply, err error)
	SayCreateCont(ctx context.Context, req *CreateContRequest, opts ...http.CallOption) (rsp *CreateContReply, err error)
	SayCreateTransaction(ctx context.Context, req *CreateTransactionRequest, opts ...http.CallOption) (rsp *CreateTransactionReply, err error)
	SayMiners(ctx context.Context, req *GetMinersRequest, opts ...http.CallOption) (rsp *GetMinersReply, err error)
	SayNodes(ctx context.Context, req *GetNodesRequest, opts ...http.CallOption) (rsp *GetNodesReply, err error)
	SayTransaction(ctx context.Context, req *GetTransactionRequest, opts ...http.CallOption) (rsp *GetTransactionReply, err error)
}

type GreeterHTTPClientImpl struct {
	cc *http.Client
}

func NewGreeterHTTPClient(client *http.Client) GreeterHTTPClient {
	return &GreeterHTTPClientImpl{client}
}

func (c *GreeterHTTPClientImpl) SayAccount(ctx context.Context, in *GetAccountRequest, opts ...http.CallOption) (*GetAccountReply, error) {
	var out GetAccountReply
	pattern := "/account/{addr}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/pro.v1.Greeter/SayAccount"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GreeterHTTPClientImpl) SayBlock(ctx context.Context, in *BlockRequest, opts ...http.CallOption) (*BlockReply, error) {
	var out BlockReply
	pattern := "/block/{parm}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/pro.v1.Greeter/SayBlock"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GreeterHTTPClientImpl) SayBlocks(ctx context.Context, in *BlocksRequest, opts ...http.CallOption) (*BlocksReply, error) {
	var out BlocksReply
	pattern := "/blocks"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/pro.v1.Greeter/SayBlocks"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GreeterHTTPClientImpl) SayCont(ctx context.Context, in *GetContRequest, opts ...http.CallOption) (*GetContReply, error) {
	var out GetContReply
	pattern := "/cont/{tx}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/pro.v1.Greeter/SayCont"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GreeterHTTPClientImpl) SayCreateCont(ctx context.Context, in *CreateContRequest, opts ...http.CallOption) (*CreateContReply, error) {
	var out CreateContReply
	pattern := "/cont"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/pro.v1.Greeter/SayCreateCont"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GreeterHTTPClientImpl) SayCreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...http.CallOption) (*CreateTransactionReply, error) {
	var out CreateTransactionReply
	pattern := "/transaction"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/pro.v1.Greeter/SayCreateTransaction"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GreeterHTTPClientImpl) SayMiners(ctx context.Context, in *GetMinersRequest, opts ...http.CallOption) (*GetMinersReply, error) {
	var out GetMinersReply
	pattern := "/miners"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/pro.v1.Greeter/SayMiners"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GreeterHTTPClientImpl) SayNodes(ctx context.Context, in *GetNodesRequest, opts ...http.CallOption) (*GetNodesReply, error) {
	var out GetNodesReply
	pattern := "/nodes"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/pro.v1.Greeter/SayNodes"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GreeterHTTPClientImpl) SayTransaction(ctx context.Context, in *GetTransactionRequest, opts ...http.CallOption) (*GetTransactionReply, error) {
	var out GetTransactionReply
	pattern := "/transaction/{tx}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/pro.v1.Greeter/SayTransaction"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}