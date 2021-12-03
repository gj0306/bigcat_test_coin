package network

import (
	"encoding/json"
	"github.com/google/wire"
	"github.com/op/go-logging"
	"go.uber.org/zap"
	"io/ioutil"
	"math/rand"
	"net/http"
	"bigcat_test_coin/app/blockchain/internal/config"
	"bigcat_test_coin/tools"
	"sync"
	"time"

	//"github.com/gogo/protobuf/proto"
	"github.com/nknorg/nnet"
	"github.com/nknorg/nnet/node"
	"github.com/nknorg/nnet/protobuf"
	"github.com/nknorg/nnet/util"
)

const dataMaxLen int = 1024
var ProviderSet = wire.NewSet(NewNetWork)

type NetWork struct {
	nn              *nnet.NNet
	log *zap.Logger
	port     uint16
	IsCreate bool
	Addr     string
	Data 			chan Message
	nodeAddress     map[string]struct{}
	conf            *config.NetworkConf
	//黑名单
	blacklist [][]byte
	mu sync.Mutex
}

func create(port uint16) (*nnet.NNet, error) {
	id, _ := util.RandBytes(32)
	conf := &nnet.Config{
		Port:                  port,
		Transport:             "kcp",
		BaseStabilizeInterval: 500 * time.Millisecond,
		NumFingerSuccessors:   2,
	}
	//日志等级设置
	logging.SetLevel(logging.ERROR, "nnet")
	nn, err := nnet.NewNNet(id, conf)
	if err != nil {
		return nil, err
	}
	return nn, nil
}

func NewNetWork(conf *config.NetworkConf,log *zap.Logger)(netWork *NetWork){
	nn, err := create(conf.Port)
	if err != nil{
		return nil
	}
	netWork = &NetWork{
		nn:          nn,
		port:        conf.Port,
		Data:        make(chan Message, dataMaxLen),
		IsCreate:    conf.IsCreate,
		Addr:        nn.GetLocalNode().Addr,
		blacklist:   make([][]byte,0),
		nodeAddress: make(map[string]struct{}),
		log:         log,
		conf:        conf,
	}
	//接受消息
	nn.MustApplyMiddleware(node.BytesReceived{Func: func(msg, msgID, srcID []byte, remoteNode *node.RemoteNode) ([]byte, bool) {
		if remoteNode != nil {
			remoteNode.GetId()
			if len(netWork.Data) < dataMaxLen && len(msg)>0{
				netWork.Data <- Message{
					Rid: remoteNode.GetId(),
					Data: msg,
				}
				netWork.addNodeAddr(remoteNode.Addr)
			}
		}
		return msg, true
	}})
	//黑名单处理
	//nn.MustApplyMiddleware(node.WillConnectToNode{
	//	Func: func(p *protobuf.Node) (bool, bool) {
	//		id := p.GetId()
	//		for _,addr := range netWork.blacklist{
	//			if tools.BytesEqual(id,addr){
	//				return false,false
	//			}
	//		}
	//		return true,true
	//	},
	//	Priority: 0,
	//})
	//断开连接
	nn.MustApplyMiddleware(node.RemoteNodeDisconnected{
		Func: func(remoteNode *node.RemoteNode) bool {
			nodes,_ := netWork.getRemoteNodes()
			if len(nodes)==0{
				netWork.Join()
			}
			return true
		},
		Priority: 0,
	})

	err = nn.Start(conf.IsCreate)
	if err != nil {
		return nil
	}
	return netWork
}

// InitWork 网络初始化操作
func (work *NetWork)InitWork(){
	if !work.IsCreate {
		work.Join()
	}
}

// Join 连接网络服务
func (work *NetWork)Join(){
	for {
		nodes,err := work.getHtos(work.conf.Addr)
		if err != nil || len(nodes)==0{
			work.log.Error("work Join 获取区块网络地址 err:"+err.Error())
			time.Sleep(time.Second*3)
			continue
		}
		err = work.nn.Join(nodes[rand.Intn(len(nodes))])
		if err == nil{
			break
		}else {
			work.log.Error("work Join 网络连接失败 err:"+err.Error())
		}
		time.Sleep(time.Second*5)
	}
	return
}
// GetAddr 获取节点地址
func (work *NetWork) GetAddr()string{
	return work.Addr
}
func (work *NetWork) SendMessage(msg []byte) (err error) {
	_, err = work.nn.SendBytesBroadcastAsync(
		msg,
		protobuf.BROADCAST_TREE,
	)
	if err != nil{
		work.log.Error("work send message",
			zap.String("err",err.Error()),
		)
	}
	return err
}
func (work *NetWork) SendOneMessage(msg []byte, destID []byte) (err error) {
	_, err = work.nn.SendBytesRelayAsync(msg, destID)
	if err != nil{
		work.log.Error("work send one message",
			zap.String("err",err.Error()),
		)
	}
	return err
}
func (work *NetWork) getRemoteNodes() (remoteNodes []*node.RemoteNode,err error) {
	localNode := work.nn.GetLocalNode()
	if localNode != nil {
		remoteNodes, err = localNode.GetNeighbors(func(remoteNode *node.RemoteNode) bool {
			//筛选条件
			return true
		})
		return remoteNodes,err
	}
	return nil,nil
}
func (work *NetWork) addNodeAddr(addr string){
	work.mu.Lock()
	work.nodeAddress[addr]= struct{}{}
	work.mu.Unlock()
}
func (work *NetWork) AddBlackList(rid []byte){
	nodes,err := work.getRemoteNodes()
	if err == nil{
		for _, remoteNode := range nodes{
			if tools.BytesEqual(remoteNode.GetId(),rid){
				_ = remoteNode.NotifyStop()
				break
			}
		}
	}
	work.blacklist = append(work.blacklist, rid)
}

func (work *NetWork) GetNodes() (nodes [][]byte,err error) {
	rs,err := work.getRemoteNodes()
	if err != nil{
		return nil, err
	}
	nodes = make([][]byte,0,len(rs))
	for _,r := range rs{
		nodes = append(nodes, r.GetId())
	}
	return nodes,nil
}
func (work *NetWork) GetNodeAddress()[]string{
	address := make([]string,0)
	address = append(address, work.GetAddr())
	rs,_ := work.getRemoteNodes()
	for _,r := range rs{
		address = append(address, r.GetAddr())
	}
	return address
}

func (work *NetWork) getHtos(addr string)(hs []string,err error){
	resp,err := http.Get("addr")
	if err != nil{
		return nil,err
	}
	body,err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil{
		return nil,err
	}
	mp := map[string][]string{}
	err = json.Unmarshal(body, &mp)
	if err != nil{
		return nil,err
	}
	return mp["nodes"],nil
}