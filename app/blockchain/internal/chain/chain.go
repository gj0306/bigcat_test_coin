package chain

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"bigcat_test_coin/app/blockchain/internal/config"

	v1 "bigcat_test_coin/app/api"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/box"
	"bigcat_test_coin/app/blockchain/internal/common"
	"bigcat_test_coin/app/blockchain/internal/network"
	"bigcat_test_coin/app/blockchain/internal/service"
	"bigcat_test_coin/tools"
	"strconv"
	"sync"
	"time"
)

type Chain struct {
	privateKey *blc.Address
	work       *network.NetWork
	cBox       *box.Box
	//配置文件
	conf *config.YamlConfig
	log *zap.Logger
	//接受数据方法
	recvApp map[byte]func(rid, data []byte) error
	//同步时间 单位毫秒
	syncId []byte
	//web
	server *service.GreeterService

	nodeInfos  map[string]*v1.P2NGetInfo
	nodeHeight int64
	nodeId     []byte
	syncTm     int64
	mu         sync.Mutex
	grpcServer *grpc.Server
	httpServer *http.Server
	ctx context.Context
}

func NewChain(cBox *box.Box,conf *config.YamlConfig, work *network.NetWork, privateKey *blc.Address, log *zap.Logger, server *service.GreeterService,grpcServer *grpc.Server,httpServer *http.Server) *Chain {
	c := &Chain{
		privateKey: privateKey,
		conf:       conf,
		work:       work,
		cBox:       cBox,
		recvApp:    map[byte]func(rid []byte, data []byte) error{},
		log:        log,
		syncId:     nil,
		nodeInfos:  nil,
		server:     server,
		mu:         sync.Mutex{},
		grpcServer: grpcServer,
		httpServer: httpServer,
		ctx: context.Background(),
	}
	//注册 通讯方法
	N2PGetInfo(c)
	P2NGetInfo(c)
	N2PGetBlock(c)
	N2PGetBlocksHash(c)
	P2NGetBlock(c)
	PlushVerifier(c)
	PushContracts(c)
	PushTransaction(c)
	return c
}

// Init 初始化
func (c *Chain) Init() {
	//数据初始化
	c.cBox.InitDb()
	//server初始化
	if c.grpcServer != nil && c.conf.Server.Open{
		go func() {
			err := c.grpcServer.Start(c.ctx)
			if err != nil{
				c.log.Error("grpc server err:" + err.Error())
			}
		}()
	}
	if c.httpServer != nil && c.conf.Server.Open{
		go func() {
			err := c.httpServer.Start(c.ctx)
			if err != nil{
				c.log.Error("http server err:" + err.Error())
			}
		}()
	}
	time.Sleep(time.Second)
	c.serviceBox()
	//网络初始化
	c.work.InitWork()
	//开始接受消息
	go c.recv()
}

// BlockSync 同步Block
func (c *Chain) BlockSync() {
	var nodes [][]byte
start:
	nodes, _ = c.work.GetNodes()
	c.log.Debug("chain BlockSync 连接节点数量" + strconv.Itoa(len(nodes)))
	time.Sleep(time.Second * 2)
	//
	c.sendInfoMessage()
	time.Sleep(time.Second * 3)
	if c.nodeHeight == 0 || len(nodes) < 1 {
		goto start
	}
	for c.nodeHeight > c.cBox.RightHeight+int64(common.BlockChainConsensusMechanismSize) {
		c.sendInfoMessage()
		time.Sleep(time.Millisecond * 500)
		continue
	}
}

func (c *Chain) Run() {
	if !c.work.IsCreate {
		c.BlockSync()
	}
	for {
		second := time.Now().Second()
		lastBlock := c.cBox.GetBlock(c.cBox.RightHeight)
		tm := common.GetNextBlockTime(lastBlock.TimeStamp, 1)
		if tm <= time.Now().Unix() {
			c.PushBlock()
			continue
		}
		if second < 5 || second > 55 {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		//同步
		c.sendInfoMessage()
		time.Sleep(time.Second * 1)
		//心跳测试
		c.Heartbeat()
	}
}

func (c *Chain) createBlock() *blc.Block {
	preBlock := c.cBox.GetBlock(c.cBox.RightHeight)
	vs := c.cBox.GetVerifiers(preBlock.Hash)
	conts, trans := c.cBox.GetContTranCache()
	block := &blc.Block{
		Contracts:    conts,
		Transactions: trans,
		Verifiers:    vs,
	}
	accs := c.cBox.GenerateAccounts(block)
	var tm int64
	if preBlock.TimeStamp > 0 {
		tm = common.GetNextBlockTime(time.Now().Unix(), 1)
	} else {
		tm = common.GetNextBlockTime(time.Now().Unix(), 0)
	}
	newBlock := blc.NewBlock(preBlock.Hash, preBlock.Height+1, trans,
		vs, int64(len(vs))+preBlock.VerifierTotal, conts, accs, tm)
	return newBlock
}
func (c *Chain) PushBlock() {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
	}()
	block := c.createBlock()
	err := c.cBox.AddBlocks([]*blc.Block{block})
	if err != nil {
		c.log.Error("push block 发块儿失败",
			zap.String("err", err.Error()),
			zap.Int64("height", block.Height),
			zap.Int64("verifierTotal", block.VerifierTotal),
			zap.Int64("tm", block.TimeStamp),
			zap.String("hash", tools.Encodeb58(block.Hash)),
		)
		return
	}
	c.cBox.Next()
	c.log.Debug("push block",
		zap.Int64("height", block.Height),
		zap.Int64("verifierTotal", block.VerifierTotal),
		zap.Int64("tm", block.TimeStamp),
		zap.String("hash", tools.Encodeb58(block.Hash)),
	)
	//service
	c.serviceBox()
	//校验信息
	if c.privateKey != nil {
		if c.cBox.IsWork(c.privateKey.GetAddress(), block.Hash) {
			v := blc.NewVerifier(c.privateKey, block.Hash)
			reply := &v1.PlushVerifier{
				Verifier: v,
			}
			bys, err := proto.Marshal(reply)
			if err != nil {
				c.log.Error("proto marshal verifier",
					zap.Int64("height", block.Height),
					zap.Int64("verifierTotal", block.VerifierTotal),
					zap.Int64("tm", block.TimeStamp),
					zap.String("hash", tools.Encodeb58(block.Hash)),
					zap.String("err", err.Error()),
				)
				return
			}
			err = c.work.SendMessage(v1.EncodingMessage(v1.PlushVerifierNumber, bys))
			c.cBox.AddVerifier(v)
			if err == nil {
				//c.cBox.AddVerifier(v)
			} else {
				c.log.Error("push verifier",
					zap.String("preHash", tools.Encodeb58(v.PreHash)),
					zap.String("publicKey", tools.Encodeb58(v.PublicKey)),
					zap.String("signature", tools.Encodeb58(v.Signature)),
					zap.String("hash", tools.Encodeb58(block.Hash)),
					zap.String("err", err.Error()),
				)
			}
		}
	}
}

func (c *Chain) sendInfoMessage() {
	//是否同步
	if c.nodeId != nil && (c.syncTm > time.Now().Unix() && c.syncTm > 0) {
		return
	}
	nodes, err := c.work.GetNodes()
	if err != nil || len(nodes)==0 {
		if err != nil{
			c.log.Error("chain sendInfoMessage err:"+err.Error())
		}
		return
	}
	c.mu.Lock()
	c.nodeInfos = make(map[string]*v1.P2NGetInfo)
	leftBlock := c.cBox.GetBlock(c.cBox.LeftHeight)
	lastBlock := c.cBox.GetBlock(c.cBox.RightHeight)
	c.mu.Unlock()
	for _, nid := range nodes {
		resp := &v1.N2PGetInfo{
			SearchHeight: leftBlock.Height,
		}
		bys, err := proto.Marshal(resp)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		err = c.work.SendOneMessage(v1.EncodingMessage(v1.N2PGetInfoNumber, bys), nid)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	<-c.echoNodeInfos(ctx, len(nodes))
	type record struct {
		number int
		reply  *v1.P2NGetInfo
		ids    [][]byte
	}
	mp := make(map[string]record)
	for key, val := range c.nodeInfos {
		if !tools.BytesEqual(val.SearchHash, leftBlock.Hash) && leftBlock.Height > int64(common.BlockChainConsensusMechanismSize) {
			continue
		}
		//权重对比
		if !common.BlockWeightRatio(&blc.Block{
			Hash:          val.LastHash,
			VerifierTotal: val.VerifierTotal,
		},lastBlock){
			continue
		}
		lastKey := string(val.LastHash)
		r, ok := mp[lastKey]
		if !ok {
			r = record{
				number: 0,
				reply:  val,
				ids:    make([][]byte, 0),
			}
		}
		r.number++
		r.ids = append(r.ids, []byte(key))
		mp[lastKey] = r
	}
	//
	var number int
	var key string
	for lastKey, r := range mp {
		if len(r.ids) > number {
			number = len(r.ids)
			key = lastKey
		}
	}
	//get block
	r := mp[key]
	if len(r.ids)<1{
		return
	}
	c.nodeHeight = r.reply.GetHeight()
	nid := r.ids[rand.Intn(len(r.ids))]
	c.syncTm = time.Now().Unix() + 3
	c.nodeId = nid
	reply := &v1.N2PGetBlock{
		LeftHeight:  leftBlock.Height,
		RightHeight: lastBlock.Height,
		LastHash:    lastBlock.Hash,
	}
	bys, err := proto.Marshal(reply)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = c.work.SendOneMessage(v1.EncodingMessage(v1.N2PGetBlockNumber, bys), nid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c *Chain) echoNodeInfos(ctx context.Context, number int) chan int {
	ch := make(chan int, 1)
	sign := true
	for sign {
		select {
		case _ = <-ctx.Done():
			sign = false
		default:
			if len(c.nodeInfos) == number {
				sign = false
			} else {
				time.Sleep(time.Millisecond * 50)
			}
		}
	}
	ch <- len(c.nodeInfos)
	return ch
}

// Heartbeat 网络心跳
func (c *Chain) Heartbeat() {
	_, err := c.work.GetNodes()
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Println(len(node))
}

// WorkJoin 网络连接
func (c *Chain) WorkJoin() {
	c.work.Join()
}

//server box
func (c *Chain) serviceBox(){
	if c.server != nil{
		c.server.WebBox = c.cBox.NewSonBox(c.cBox.RightHeight)
	}
}

//接受消息
func (c *Chain) recv() {
	for msg := range c.work.Data {
		f, ok := c.recvApp[msg.Data[0]]
		if ok {
			c.mu.Lock()
			err := f(msg.Rid, msg.Data[1:])
			c.mu.Unlock()
			if err == nil {
				continue
			}
			fmt.Println(err.Error())
		}
	}
}
