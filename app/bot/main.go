package main

import (
	"context"
	"fmt"
	"math/rand"
	v1 "bigcat_test_coin/app/api"
	"bigcat_test_coin/app/bot/bot"
	"bigcat_test_coin/app/bot/client"
	"time"
)

var botAddrs = []string{
	"9d7uj6PkeHMxwiLv5fzzZQWRuSuXXToLyjf3F1NGdrtv",
	"1AAMT3FKnjHvdJjT1enQrD6DSW6LqWADYwLionLuWFeq",
	"DkMDVi61vhFPB3rpM3yQbM3n1tLxF4f1NZphnLmUm5YQ",
	"HMhjAwFsrNLDKfatvi4ftNvx4Lx5Xrywt8XMdGG4NTyV",
	"Gx1mouB1PRngihP6P5tRoNfzVatriCE13h7tdgbzKNLe",
}
var bots = []*bot.Bot{}
const endpoint = "127.0.0.1:9902"

func main()  {
	var count int64
	grpcClient := client.NewOrderServiceClient(endpoint)
	for _,addr := range botAddrs{
		bots = append(bots, bot.NewBot(addr))
	}
	n := len(bots)
	ctx := context.Background()
	for {
		form := bots[rand.Intn(n)]
		to  := bots[rand.Intn(n)]
		if form.Addr == to.Addr{
			continue
		}
		switch rand.Intn(10) {
		case 0:
			tran := form.NewMiner()
			_,err := grpcClient.SayCreateTransaction(ctx,&v1.CreateTransactionRequest{Tran: tran})
			if err != nil{
				fmt.Println(err.Error())
			}
		case 1,2,3,4,5:
			v,err := grpcClient.SayAccount(ctx,&v1.GetAccountRequest{Addr: form.Addr})
			if err != nil{
				fmt.Println(err.Error())
				break
			}
			if v.Account != nil{
				form.SetVal(v.Account.Value)
			}
			tran := form.NewTran(to.Addr)
			_,err = grpcClient.SayCreateTransaction(ctx,&v1.CreateTransactionRequest{Tran: tran})
			if err != nil{
				fmt.Println(err.Error())
			}
		default:
			cont := form.NewCont()
			_,err := grpcClient.SayCreateCont(ctx,&v1.CreateContRequest{Cont: cont})
			if err != nil{
				fmt.Println(err.Error())
			}
		}
		count++
		fmt.Println("计数:",count)
		time.Sleep(time.Second*2)
	}

}
