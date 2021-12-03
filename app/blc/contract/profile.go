package contract

import (
	"github.com/vmihailenco/msgpack"
	"bigcat_test_coin/app/blc"
)


func CreateProfile(name,head,introduce string) Profile {
	return Profile{
		Name:   name,
		Head:   head,
		Introduce: introduce,
	}
}

// GetType 获取类型
func (p *Profile) GetType() int64 {
	return ConnTypeProfile
}

// Serialize 序列化成[]byte
func (p *Profile) Serialize() []byte {
	b, err := msgpack.Marshal(p)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// Deserialize 反序列化
func (p *Profile) Deserialize(d []byte) error {
	err := msgpack.Unmarshal(d, p) // 将二进制流转化回结构体
	if err != nil {
		return err
	}
	return nil
}

// Check 校验
func (p *Profile) Check(cont *blc.Cont) bool {
	if len(p.Name)==0 || len(p.Head)==0|| len(p.Introduce)==0{
		return false
	}
	return true
}