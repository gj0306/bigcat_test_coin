package blc

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"bigcat_test_coin/tools"
)

// NewVerifier 生成校验
func NewVerifier(address *Address,PreHash []byte)(v *Verifier){
	signature := EllipticCurveSign(address.PrivateKey, PreHash[:])
	return &Verifier{
		Signature: signature,
		PublicKey: address.GetPublicKey(),
		PreHash:   PreHash,
	}
}

// VerifiersSort 验证信息排序
func VerifiersSort(verifiers []Verifier){
	tools.SortStructList(verifiers,[]string{"PublicKey","PreHash"})
}

// Serialize 将account序列化成[]byte
func (x *Verifier) Serialize() []byte {
	bys, err := proto.Marshal(x)
	if err != nil {
		panic(err.Error())
	}
	return bys
}

func (x *Verifier) Deserialize(d []byte) {
	err := proto.Unmarshal(d, x)
	if err != nil {
		panic(err.Error())
	}
}

func (x *Verifier)Check()bool{
	//哈希校验
	if !EllipticCurveVerify(x.PublicKey, x.Signature, x.PreHash) {
		return false
	}
	return true
}

func (x *Verifier) GetForm()string{
	return GetAddressFromPublicKey(x.PublicKey)
}

func (x *Verifier) MarshalJSON()([]byte, error){
	mp := map[string]interface{}{
		"signature":tools.Encodeb58(x.Signature),
		"publicKey":tools.Encodeb58(x.PublicKey),
		"preHash":tools.Encodeb58(x.PreHash),
	}
	return json.Marshal(mp)
}