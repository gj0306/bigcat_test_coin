package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.uber.org/zap"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/box"
	"bigcat_test_coin/app/blockchain/internal/chain"
	"bigcat_test_coin/app/blockchain/internal/coinlog"
	"bigcat_test_coin/app/blockchain/internal/config"
	"bigcat_test_coin/app/blockchain/internal/network"
	"bigcat_test_coin/app/blockchain/internal/service"
	"strconv"
)

func newApp(cBox *box.Box,conf *config.YamlConfig,work *network.NetWork,privateKey *blc.Address,log *zap.Logger,srv *service.GreeterService,grpcServer *grpc.Server,httpServer *http.Server) *chain.Chain {
	return chain.NewChain(cBox,conf ,work ,privateKey,log,srv,grpcServer,httpServer)
}

var (
	key  string
	c string
	initial string
	p string
	d string

	conf *config.YamlConfig
)

func init() {
	flag.StringVar(&key, "key", "9d7uj6PkeHMxwiLv5fzzZQWRuSuXXToLyjf3F1NGdrtv", "")
	flag.StringVar(&c, "conf", "./app/blockchain/configs/config.yaml", "config path, eg: -conf config.yaml")

	flag.StringVar(&initial, "initial", "ok", "")
	flag.StringVar(&p, "port", "", "")
	flag.StringVar(&d, "data", "", "")
}

func main()  {
	flag.Parse()
	//配置
	conf = config.NewConfig(c)
	if conf == nil{
		panic("配置文件加载失败")
	}
	if initial!=""&&initial!="false"{
		conf.Network.IsCreate = true
	}
	if p !=""{
		port,err := strconv.ParseUint(p,10,64)
		if err == nil{
			conf.Network.Port = uint16(port)
		}
	}
	//数据库地址
	if d != ""{
		conf.Db.Dir = d
	}
	//web 日志
	logger := coinlog.NewZapLog("logs","web",nil)
	//私钥
	if key == ""{
		key = conf.Address.PrivateKey
	}
	privateKey,err := blc.LoadAddress(key)
	if err != nil{
		panic(err)
	}
	//cil
	cil := NewCli()
	cil.Run()
	//注册生成服务
	app, cleanup, err := initApp(privateKey,conf,conf.Log,conf.Db,conf.Address,conf.Network,conf.Server,logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	app.Init()
	app.Run()
}



