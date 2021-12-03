package blc

import (
	"crypto/sha256"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"bigcat_test_coin/tools"
)


// connInterface 合约接口
type connInterface interface {
	GetType() int64
	Serialize() []byte
	Deserialize(d []byte)error
	Check(cont *Cont) bool
}


func NewConn(c connInterface,privateKey string) (*Cont,error){
	conn := &Cont{
		ConnType:  c.GetType(),
		Data:      c.Serialize(),
		Number:    rand.Int63(),
	}
	paddr,err := LoadAddress(privateKey)
	if err != nil{
		return conn,err
	}
	conn.PublicKey=paddr.GetPublicKey()
	conn.TxHash = conn.GetHash()
	conn.Signature = EllipticCurveSign(paddr.PrivateKey,conn.TxHash)
	return conn, nil
}

func (x *Cont) Serialize() []byte {
	bys, err := proto.Marshal(x)
	if err != nil {
		panic(err.Error())
	}
	return bys
}

func (x *Cont) Deserialize(d []byte) {
	err := proto.Unmarshal(d, x)
	if err != nil {
		panic(err.Error())
	}
}

// GetHash 获取哈希
func (x *Cont) GetHash()[]byte{
	conn := &Cont{
		ConnType:  x.ConnType,
		PublicKey: x.PublicKey,
		Data:      x.Data,
		Quote:     x.Quote,
		Number:    x.Number,
	}
	hash := sha256.Sum256(conn.Serialize())
	return hash[:]
}

// GetForm 获取发起人地址
func (x *Cont) GetForm()string {
	return GetAddressFromPublicKey(x.PublicKey)
}

// Check 校验
func (x *Cont) Check()bool  {
	//哈希校验
	hash := x.GetHash()
	if !tools.BytesEqual(hash, x.TxHash){
		return false
	}
	//地址校验
	if !EllipticCurveVerify(x.PublicKey, x.Signature, x.TxHash){
		return false
	}
	return true
}

func (x *Cont)MarshalJSON()([]byte, error){
	mp := map[string]interface{}{
		"connType":x.ConnType,
		"publicKey":tools.Encodeb58(x.PublicKey),
		"signature":tools.Encodeb58(x.Signature),
		"txHash":tools.Encodeb58(x.TxHash),
		"number":x.Number,
		"data":x.Data,
	}
	return json.Marshal(mp)
}