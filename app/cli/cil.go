package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	v1 "bigcat_test_coin/app/api"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/cli/client"
	"os"
	"strconv"
	"strings"
)

var private *blc.Address

type Cli struct {
}

//打印帮助提示
func printUsage() {
	fmt.Println("----------------------------------------------------------------------------- ")
	fmt.Println("Usage:")
	fmt.Println("\thelp                                             打印命令行说明")
	fmt.Println("\tsetServer  -s DATA[服务器地址]                     设置当前cil连接服务器")
	fmt.Println("\tsetWallets -p DATA[私钥]                          设置当前cil钱包")
	fmt.Println("\tbalance                                           查看当前钱包账户余额")
	fmt.Println("\tbalance -a DATA[地址]                             查看指定账户余额")
	fmt.Println("\ttransfer -to DATA[对方地址] -amount DATA[数量]      进行转账操作")
	fmt.Println("\tblock -h data[高度] -tx data[哈希]                 查看Block信息")
	fmt.Println("\tlastBlocks                                        查看最后发布的blocks")
	fmt.Println("\ttransaction -tx DATA[哈希]                         查看交易")
	fmt.Println("\tcontract -tx DATA[哈希]                            查看合约")

	fmt.Println("------------------------------------------------------------------------------")
}

func New() *Cli {
	return &Cli{}
}

func (cli *Cli) Run() {
	printUsage()
	//go cli.startNode()
	cli.ReceiveCMD()
}

// ReceiveCMD 获取用户输入
func (cli Cli) ReceiveCMD() {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}
		cli.userCmdHandle(sendData)
	}
}

//用户输入命令的解析
func (cli Cli) userCmdHandle(data string) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(fmt.Println(err))
		}
	}()
	ctx := context.Background()
	//去除命令前后空格
	data = strings.TrimSpace(data)
	var cmd string
	var msg string
	if strings.Contains(data, " ") {
		cmd = data[:strings.Index(data, " ")]
		msg = data[strings.Index(data, " ")+1:]
	} else {
		cmd = data
	}
	switch cmd {
	case "help":
		printUsage()
		fmt.Println(msg)
	case "setServer":
		host := getSpecifiedContent(data, "-s")
		client.SetClientEndpoint(host)
		c := client.GetClient()
		_,err := c.SayBlock(ctx,&v1.BlockRequest{Parm: strconv.Itoa(1)})
		if err != nil{
			fmt.Println(err.Error())
		}else {
			fmt.Println("设置成功")
		}
	case "setWallets":
		address := getSpecifiedContent(data, "-p")
		pr,err := blc.LoadAddress(address)
		if err != nil{
			fmt.Println(err.Error())
		}else {
			private = pr
			fmt.Println("设置钱包成功")
		}
	case "balance":
		addr := getSpecifiedContent(data, "-a")
		if addr == ""{
			if private == nil{
				fmt.Println("尚未设置钱包地址")
				break
			}
			addr = private.GetAddress()
		}
		c := client.GetClient()
		rep,err := c.SayAccount(ctx,&v1.GetAccountRequest{Addr:addr})
		if err != nil{
			fmt.Println(err.Error())
		}else {
			if rep.Account == nil{
				fmt.Println("数据不存在")
			}else {
				bys,_ := json.Marshal(rep.Account)
				fmt.Println(string(bys))
			}
		}
	case "transfer":
		if private == nil{
			if private == nil{
				fmt.Println("尚未设置钱包地址")
				break
			}
		}
		to := getSpecifiedContent(data, "-to")
		amount := getSpecifiedContent(data, "-amount")
		val,err := strconv.ParseInt(amount,10,64)
		if err != nil{
			fmt.Println(err.Error())
			break
		}
		if !blc.IsVerifyAddress(to){
			fmt.Println("非有效地址")
			break
		}
		c := client.GetClient()
		rep,err := c.SayAccount(ctx,&v1.GetAccountRequest{Addr: private.GetAddress()})
		if err != nil{
			fmt.Println(err.Error())
			break
		}
		if rep.Account == nil{
			fmt.Println("余额不足")
			break
		}
		if rep.Account.Value < val{
			fmt.Println("余额不足")
			break
		}
		tran := &blc.Transaction{
			To:        to,
			Value:     val,
			Fee:       0,
			Number:    rand.Int63(),
			PublicKey: private.GetPublicKey(),
			Signature: nil,
			TxHash:    nil,
		}
		tran.TxHash = tran.GetHash()
		tran.Signature = blc.EllipticCurveSign(private.PrivateKey,tran.TxHash)
		_,err = c.SayCreateTransaction(ctx,&v1.CreateTransactionRequest{Tran: tran})
		if err != nil{
			fmt.Println(err.Error())
			break
		}
		fmt.Println("发起转账合约成功")
	case "block":
		var parm string
		h := getSpecifiedContent(data, "-h")
		if h != ""{
			_,err := strconv.ParseInt(h,10,64)
			if err != nil{
				fmt.Println(err.Error())
				break
			}
			parm = h
		}
		tx := getSpecifiedContent(data, "-tx")
		if tx != ""{
			parm = tx
		}
		if parm == ""{
			fmt.Println("缺少必要参数")
			break
		}
		c := client.GetClient()
		rep,err := c.SayBlock(ctx,&v1.BlockRequest{Parm:parm})
		if err != nil{
			fmt.Println(err.Error())
			break
		}
		if rep.Block != nil{
			bys,_ := json.Marshal(rep.Block)
			fmt.Println(string(bys))
		}else {
			fmt.Println("未找到block")
		}
	case "lastBlocks":
		c := client.GetClient()
		rep,err := c.SayBlocks(ctx,&v1.BlocksRequest{})
		if err != nil{
			fmt.Println(err.Error())
			break
		}
		for _,block := range rep.Blocks{
			bys,_ := json.Marshal(block)
			fmt.Println(string(bys))
		}
	case "transaction":
		tx := getSpecifiedContent(data, "-tx")
		if tx == ""{
			fmt.Println("缺少必要参数")
			break
		}
		c := client.GetClient()
		rep,err := c.SayTransaction(ctx,&v1.GetTransactionRequest{Tx: tx})
		if err != nil{
			fmt.Println(err.Error())
			break
		}
		if rep.Transaction == nil{
			fmt.Println("未找到 Transaction")
		}else {
			bys,_ := json.Marshal(rep.Transaction)
			fmt.Println(string(bys))
		}
	case "contract":
		tx := getSpecifiedContent(data, "-tx")
		if tx == ""{
			fmt.Println("缺少必要参数")
			break
		}
		c := client.GetClient()
		rep,err := c.SayCont(ctx,&v1.GetContRequest{Tx: tx})
		if err != nil{
			fmt.Println(err.Error())
			break
		}
		if rep.Cont == nil{
			fmt.Println("未找到 Contract")
		}else {
			bys,_ := json.Marshal(rep.Cont)
			fmt.Println(string(bys))
		}
	default:
		fmt.Println("无此命令!")
		printUsage()
	}
}

//返回data字符串中,标签为tag的内容
func getSpecifiedContent(data, tag string) string {
	sign := false
	for _,v := range strings.Split(data," "){
		if sign{
			if strings.Index(v,"-") == 0{
				return ""
			}
			return v
		}
		if v == tag{
			sign = true
		}
	}
	return ""
}
