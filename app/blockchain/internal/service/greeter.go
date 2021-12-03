package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/proto"
	v1 "bigcat_test_coin/app/api"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/box"
	"bigcat_test_coin/app/blockchain/internal/network"
	"bigcat_test_coin/tools"
	"strconv"
	"sync"
	"time"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer
	log    *log.Helper
	WebBox *box.Box
	Box    *box.Box
	work   *network.NetWork
	_trans []*blc.Transaction
	_conts []*blc.Cont
	mu      sync.Mutex
}

// NewGreeterService new a greeter service.
func NewGreeterService(logger log.Logger,box *box.Box,work *network.NetWork) *GreeterService {
	return &GreeterService{ log: log.NewHelper(logger),Box: box,work:work}
}

// SayBlocks implements api.SayBlocks
func (s *GreeterService) SayBlocks(ctx context.Context, in *v1.BlocksRequest) (*v1.BlocksReply, error) {
	left  := in.GetLeftHeight()
	right := in.GetRightHeight()
	if left == 0 || right == 0{
		left = s.WebBox.LeftHeight
		right = s.WebBox.RightHeight
	}
	reply := &v1.BlocksReply{Blocks:s.WebBox.GetBlocks(left,right)}
	return reply,nil
}

func (s *GreeterService) SayBlock(ctx context.Context, in *v1.BlockRequest) (*v1.BlockReply, error) {
	reply := &v1.BlockReply{}
	height,err := strconv.ParseInt(in.Parm,10,64)
	switch err {
	case nil:
		reply.Block = s.WebBox.GetBlock(height)
	default:
		reply.Block = s.WebBox.GetBlockByHash(tools.Decodeb58(in.Parm))
	}
	return reply,nil
}

func (s *GreeterService) SayTransaction(ctx context.Context, in *v1.GetTransactionRequest) (*v1.GetTransactionReply, error) {
	reply := &v1.GetTransactionReply{
		Transaction: s.WebBox.GetTransaction(tools.Decodeb58(in.GetTx())),
	}
	return reply,nil
}

func (s *GreeterService) SayCont(ctx context.Context, in *v1.GetContRequest) (*v1.GetContReply, error) {
	reply := &v1.GetContReply{
		Cont: s.WebBox.GetContract(tools.Decodeb58(in.GetTx())),
	}
	return reply,nil
}

func (s *GreeterService) SayAccount(ctx context.Context, in *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
	reply := &v1.GetAccountReply{
		Account: s.WebBox.GetAccount(in.GetAddr()),
	}
	return reply,nil
}

func (s *GreeterService) SayMiners(ctx context.Context, in *v1.GetMinersRequest) (*v1.GetMinersReply, error) {
	reply := &v1.GetMinersReply{
		Miners: s.WebBox.GetMiners(),
	}
	return reply,nil
}

func (s *GreeterService) SayNodes(ctx context.Context, in *v1.GetNodesRequest) (*v1.GetNodesReply, error) {
	s.work.GetAddr()
	reply := &v1.GetNodesReply{
		Nodes: s.work.GetNodeAddress(),
	}
	return reply,nil
}

func (s *GreeterService) SayCreateTransaction(ctx context.Context, in *v1.CreateTransactionRequest) (*v1.CreateTransactionReply, error){
	var err error
	if in.Tran != nil{
		second := time.Now().Second()
		if second>5 && second<=40{
			reply := &v1.PushTransactions{
				Trans: []*blc.Transaction{in.Tran},
			}
			bys, _ := proto.Marshal(reply)
			err = s.work.SendMessage(v1.EncodingMessage(v1.PushTransactionsNumber, bys))
			if err != nil{
				s.log.Error("GreeterService SayCreateTransaction err:",err.Error())
			}else {
				s.Box.AddTransaction(in.Tran)
			}
		}else {
			s.mu.Lock()
			switch s._trans {
			case nil:
				s._trans = make([]*blc.Transaction,0)
				s._trans = append(s._trans, in.Tran)
				go func() {
					time.Sleep(time.Duration(70-second)*time.Second)
					reply := &v1.PushTransactions{
						Trans: s._trans,
					}
					bys, _ := proto.Marshal(reply)
					err = s.work.SendMessage(v1.EncodingMessage(v1.PushTransactionsNumber, bys))
					if err != nil{
						s.log.Error("GreeterService SayCreateTransaction err:",err.Error())
					}else {
						for _,tran := range s._trans{
							s.Box.AddTransaction(tran)
						}
					}
					s._trans = nil
				}()
			default:
				sign := true
				for _,tran := range s._trans{
					if tools.BytesEqual(tran.TxHash,in.Tran.TxHash){
						sign = false
						break
					}
				}
				if sign{
					s._trans = append(s._trans, in.Tran)
				}
			}
			s.mu.Unlock()
		}
	}
	msg := "ok"
	if err != nil{
		msg = err.Error()
	}
	return &v1.CreateTransactionReply{
		Code: 0,
		Msg: msg,
	},nil
}

func (s *GreeterService) SayCreateCont(ctx context.Context, in *v1.CreateContRequest) (*v1.CreateContReply, error){
	var err error
	if in.Cont != nil{
		second := time.Now().Second()
		if second>5 && second<=40{
			reply := &v1.PushContracts{
				Conts: []*blc.Cont{in.Cont},
			}
			bys, _ := proto.Marshal(reply)
			err = s.work.SendMessage(v1.EncodingMessage(v1.PushContractsNumber, bys))
			if err != nil{
				s.log.Error("GreeterService SayCreateCont err:",err.Error())
			}else {
				s.Box.AddContract(in.Cont)
			}
		}else {
			s.mu.Lock()
			switch s._conts {
			case nil:
				s._conts = make([]*blc.Cont,0)
				s._conts = append(s._conts, in.Cont)
				go func() {
					time.Sleep(time.Duration(70-second)*time.Second)
					reply := &v1.PushContracts{
						Conts: s._conts,
					}
					bys, _ := proto.Marshal(reply)
					err = s.work.SendMessage(v1.EncodingMessage(v1.PushContractsNumber, bys))
					if err != nil{
						s.log.Error("GreeterService SayCreateTransaction err:",err.Error())
					}else {
						for _,conts := range s._conts{
							s.Box.AddContract(conts)
						}
					}
					s._conts = nil
				}()
			default:
				sign := true
				for _,cont := range s._conts{
					if tools.BytesEqual(cont.TxHash,in.Cont.TxHash){
						sign = false
						break
					}
				}
				if sign{
					s._conts = append(s._conts, in.Cont)
				}
			}
			s.mu.Unlock()
		}
	}
	msg := "ok"
	if err != nil{
		msg = err.Error()
	}
	return &v1.CreateContReply{
		Code: 0,
		Msg: msg,
	},nil
}