package blc

import (
	"crypto/sha256"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"bigcat_test_coin/tools"
	"time"
)

//NewBlock 生成块儿
func NewBlock(preHash []byte, height int64,ts []*Transaction, vs []*Verifier, VerifierTotal int64,cons []*Cont,as []*Account,  timeStamp int64) *Block {
	if timeStamp == 0 {
		now := time.Now()
		timeStamp = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.Local).Unix()
	}
	//排序 为了保证数据的一致性
	tools.SortStructList(vs,[]string{"PublicKey","PreHash"})
	tools.SortStructList(ts,[]string{"TxHash"})
	tools.SortStructList(cons,[]string{"Quantity","TxHash"})
	tools.SortStructList(as,[]string{"Address","Value"})

	block := &Block{
		PreHash:       preHash,
		Contracts:     cons,
		Transactions:  ts,
		TimeStamp:     timeStamp,
		Height:        height,
		Verifiers:     vs,
		VerifierTotal: VerifierTotal,
		Accounts: as,
	}
	block.Hash = block.CreateHash()
	return block
}

// Serialize 将Block对象序列化成[]byte
func (x *Block) Serialize() []byte {
	if len(x.Contracts) == 0{
		x.Contracts = nil
	}
	if len(x.Verifiers) == 0{
		x.Verifiers = nil
	}
	if len(x.Transactions) == 0{
		x.Transactions = nil
	}
	if len(x.Accounts) == 0{
		x.Accounts = nil
	}
	bys, err := proto.Marshal(x)
	if err != nil {
		panic(err.Error())
	}
	return bys
}

func (x *Block) Deserialize(d []byte) {
	err := proto.Unmarshal(d, x)
	if err != nil {
		panic(err.Error())
	}
}

func (x *Block) jsonData()(data []byte){
	obj := &Block{
		PreHash:       x.PreHash,
		Contracts:     x.Contracts,
		Transactions:  x.Transactions,
		Accounts:      x.Accounts,
		TimeStamp:     x.TimeStamp,
		Height:        x.Height,
		Verifiers:     x.Verifiers,
		VerifierTotal: x.VerifierTotal,
	}
	if len(obj.Transactions) == 0{
		obj.Transactions = nil
	}
	if len(obj.PreHash) == 0{
		obj.PreHash = nil
	}
	if len(obj.Contracts) == 0{
		obj.Contracts = nil
	}
	if len(obj.Accounts) == 0{
		obj.Accounts = nil
	}
	if len(obj.Verifiers) == 0{
		obj.Verifiers = nil
	}
	bys,_ := json.Marshal(obj)
	return bys
}

func (x *Block) CreateHash()[]byte{
	hashByte := sha256.Sum256(x.jsonData())
	return hashByte[:]
}

func (x *Block)MarshalJSON()([]byte, error){
	mp := map[string]interface{}{
		"preHash":tools.Encodeb58(x.PreHash),
		"hash":tools.Encodeb58(x.Hash),
		"timeStamp":x.TimeStamp,
		"height":x.Height,
		"verifierTotal":x.VerifierTotal,
		"contracts":x.Contracts,
		"transactions":x.Transactions,
		"accounts":x.Accounts,
		"verifiers":x.Verifiers,
	}
	return json.Marshal(mp)
}