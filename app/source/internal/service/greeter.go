package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"math/rand"
	v1 "bigcat_test_coin/app/api"
	contapi "bigcat_test_coin/app/api/contapi"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blc/contract"
	"bigcat_test_coin/tools"
)

// GreeterService is a greeter service.
type GreeterService struct {
	contapi.UnimplementedSourceGreeterServer

	log  *log.Helper
	conn v1.GreeterClient
}


// NewGreeterService new a greeter service.
func NewGreeterService(conn v1.GreeterClient,logger log.Logger,) *GreeterService {
	return &GreeterService{ log: log.NewHelper(logger),conn: conn}
}

// SaySource ...
func (s *GreeterService) SaySource(ctx context.Context, in *contapi.SourceRequest) (*contapi.SourceReply, error) {
	req,err := s.conn.SayCont(ctx,&v1.GetContRequest{
		Tx: in.Tx,
	})
	if err != nil{
		return nil, err
	}
	if req.Cont == nil{
		return &contapi.SourceReply{Source: nil},nil
	}
	if req.Cont.ConnType != contract.ConnTypeSource{
		return &contapi.SourceReply{Source: nil},nil
	}
	source := &contract.Source{}
	err = source.Deserialize(req.Cont.Data)
	if err != nil{
		return &contapi.SourceReply{Source: nil},fmt.Errorf("数据异常")
	}
	if source.Check(req.Cont){
		return &contapi.SourceReply{Source: nil},fmt.Errorf("数据异常")
	}
	return &contapi.SourceReply{Source:source},nil
}

func (s *GreeterService) SayChainSource(ctx context.Context, in *contapi.SourceChainRequest) (*contapi.SourceChainReply, error) {
	rep := &contapi.SourceChainReply{Sources: make([]*contract.Source,0)}
	tx := in.Tx
	for tx !=""{
		req,err := s.conn.SayCont(ctx,&v1.GetContRequest{
			Tx: tx,
		})
		if err != nil{
			if err != nil{
				return rep, err
			}
		}
		if req.Cont == nil{
			break
		}
		if req.Cont.ConnType != contract.ConnTypeSource{
			return rep,nil
		}
		source := &contract.Source{}
		err = source.Deserialize(req.Cont.Data)
		if err != nil{
			return rep,fmt.Errorf("合约数据异常")
		}
		if source.Check(req.Cont){
			return rep,fmt.Errorf("校验数据不通过")
		}
		rep.Sources = append(rep.Sources, source)
		if source.Quote != nil{
			tx = tools.Encodeb58(source.Quote)
		}
	}
	return rep,nil
}

func (s *GreeterService) SayCreateSource(ctx context.Context, in *contapi.CreateSourceRequest) (*contapi.CreateSourceReply, error) {
	if in.Source == nil{
		return &contapi.CreateSourceReply{
			Msg: "你真调皮！",
		},fmt.Errorf("数据为空")
	}
	addr,err := blc.LoadAddress(in.Private)
	if err != nil{
		return &contapi.CreateSourceReply{
			Msg: "私钥错误！",
		},fmt.Errorf("私钥错误")
	}

	cont := &blc.Cont{
		ConnType:  contract.ConnTypeSource,
		PublicKey: addr.GetPublicKey(),
		Data:      in.Source.Serialize(),
		Number:    rand.Int63(),
	}
	cont.TxHash = cont.GetHash()
	cont.Signature = blc.EllipticCurveSign(addr.PrivateKey,cont.TxHash)
	if !in.Source.Check(cont){
		return &contapi.CreateSourceReply{
			Msg: "校验失败！",
		},fmt.Errorf("校验失败")
	}

	return &contapi.CreateSourceReply{
		Msg: "ok！",
		Hx: tools.Encodeb58(cont.TxHash),
	},nil
}

