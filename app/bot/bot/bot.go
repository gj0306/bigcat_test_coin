package bot

import (
	"math/rand"
	"bigcat_test_coin/app/blc"
)

type Bot struct {
	Private *blc.Address
	Addr string
	Val  int64
}

func NewBot(addr string) *Bot{
	bot := &Bot{
		Private: nil,
		Addr:    "",
		Val:     0,
	}
	private,err := blc.LoadAddress(addr)
	if err != nil{
		panic(err.Error())
	}
	bot.Private = private
	bot.Addr = bot.Private.GetAddress()
	return bot
}

func (b *Bot)NewTran(to string)*blc.Transaction{
	var val int64 = 3
	if b.Val>0{
		val = rand.Int63n(b.Val)
	}
	rep := &blc.Transaction{
		To:        to,
		Value:     val,
		Fee:       0,
		Number:    rand.Int63(),
		PublicKey: b.Private.GetPublicKey(),
		Signature: nil,
		TxHash:    nil,
	}
	rep.TxHash = rep.GetHash()
	rep.Signature = blc.EllipticCurveSign(b.Private.PrivateKey,rep.TxHash)
	return rep
}

func (b *Bot)NewCont()*blc.Cont{
	bys := make([]byte, rand.Intn(4096))
	_,_ = rand.Read(bys)
	rep := &blc.Cont{
		ConnType:  rand.Int63(),
		PublicKey: b.Private.GetPublicKey(),
		Signature: nil,
		TxHash:    nil,
		Data:      bys,
		Number:    rand.Int63(),
	}
	rep.TxHash = rep.GetHash()
	rep.Signature = blc.EllipticCurveSign(b.Private.PrivateKey,rep.TxHash)
	return rep
}

func (b *Bot)NewMiner()*blc.Transaction{
	rep := &blc.Transaction{
		To:        "1M8QQfDtxJvNsE3Cw7SXA551zndDzwJT2T",
		Value:     10,
		Number:    rand.Int63(),
		PublicKey: b.Private.GetPublicKey(),
		Signature: nil,
		TxHash:    nil,
	}
	rep.TxHash = rep.GetHash()
	rep.Signature = blc.EllipticCurveSign(b.Private.PrivateKey,rep.TxHash)
	return rep
}

func (b *Bot)SetVal(val int64)  {
	b.Val = val
}
