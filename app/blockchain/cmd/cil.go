package main

import (
	"bufio"
	"fmt"
	"bigcat_test_coin/app/blc"
	"os"
	"strings"
)


type Cli struct {
}


//打印帮助提示
func printUsage() {
	fmt.Println("----------------------------------------------------------------------------- ")
	fmt.Println("Usage:")
	fmt.Println("\tstart                                             启动节点服务")
	fmt.Println("\tcreate                                            创建新的区块网络")
	fmt.Println("\tsetIp -a data[远程节点ip地址]                       设置远程节点网络ip地址")
	fmt.Println("\tsetAddr -k data[私钥地址]                          设置私钥地址")
	fmt.Println("\tnewAddr                                           生成私钥")
	fmt.Println("\tdelDb                                             删除db数据")

	fmt.Println("------------------------------------------------------------------------------")
}



func NewCli() *Cli {
	return &Cli{}
}

func (cli *Cli) Run() {
	printUsage()
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
		if !cli.userCmdHandle(sendData){
			break
		}
	}
}
func (cli Cli) userCmdHandle(data string)bool{
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(fmt.Println(err))
		}
	}()
	//去除命令前后空格
	data = strings.TrimSpace(data)
	var cmd string
	var msg string
	switch strings.Contains(data, " ") {
	case true:
		cmd = data[:strings.Index(data, " ")]
		msg = data[strings.Index(data, " ")+1:]
	case false:
		cmd = data
	}
	switch cmd {
	case "help":
		printUsage()
		fmt.Println(msg)
	case "start":
		conf.Network.IsCreate = false
		fmt.Println("开始连接区块网络")
		return false
	case "create":
		conf.Network.IsCreate = true
		fmt.Println("开始创建区块网络")
		return false
	case "setIp":
		addr := getSpecifiedContent(data, "-a")
		switch addr {
		case "":
			fmt.Println("输入错误")
		default:
			conf.Network.Addr = addr
			fmt.Println("设置节点网络地址成功")
		}

	case "setAddr":
		address := getSpecifiedContent(data, "-k")
		switch !blc.IsVerifyAddress(address) {
		case false:
			fmt.Println("地址错误")
		case true:
			key = address
			fmt.Println("设置节点地址成功")
		}
	case "newAddr":
		addr := blc.NewRandAddress()
		fmt.Println("私钥:",addr.GetPrivateKey())
		fmt.Println("地址:",addr.GetAddress())
	case "delDb":
		err := os.RemoveAll("./data")
		switch err {
		case nil:
			fmt.Println("数据删除完成")
		default:
			fmt.Println("数据删除失败 err:" + err.Error())
		}
	default:
		fmt.Println("无此命令!")
		printUsage()
	}
	return true
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