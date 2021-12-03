package blc

import (
	"crypto/sha256"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"bigcat_test_coin/tools"
)



// Serialize 将account序列化成[]byte
func (x *Transaction) Serialize() []byte {
	bys, err := proto.Marshal(x)
	if err != nil {
		panic(err.Error())
	}
	return bys
}

func (x *Transaction) Deserialize(d []byte) {
	err := proto.Unmarshal(d, x)
	if err != nil {
		panic(err.Error())
	}
}

func (x *Transaction) GetForm()string{
	return GetAddressFromPublicKey(x.PublicKey)
}

func (x *Transaction) GetHash()[]byte{
	obj := Transaction{
		To:        x.To,
		Value:     x.Value,
		Fee:       x.Fee,
		Number:    x.Number,
		PublicKey: x.PublicKey,
	}
	hashByte := sha256.Sum256(obj.Serialize())
	return hashByte[:]
}

// Check 格式校验
func (x *Transaction)Check() bool{
	//sign校验
	if !EllipticCurveVerify(x.PublicKey, x.Signature, x.TxHash) {
		return false
	}
	//地址校验
	if !IsVerifyAddress(x.GetForm()){
		return false
	}
	//哈希校验
	if !tools.BytesEqual(x.GetHash(),x.TxHash){
		return false
	}
	return true
}

func (x *Transaction)MarshalJSON()([]byte, error){
	mp := map[string]interface{}{
		"form":x.GetForm(),
		"to":x.To,
		"value":x.Value,
		"fee":x.Fee,
		"publicKey":tools.Encodeb58(x.PublicKey),
		"signature":tools.Encodeb58(x.Signature),
		"txHash":tools.Encodeb58(x.TxHash),
		"number":x.Number,
	}
	return json.Marshal(mp)
}