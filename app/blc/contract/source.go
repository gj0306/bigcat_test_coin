package contract

import (
	"github.com/vmihailenco/msgpack"
	"bigcat_test_coin/app/blc"
)

func CreateSource(from, to, content string, quote []byte) Source {
	return Source{
		From:    from,
		To:      to,
		Content: content,
		Quote:   quote,
	}
}

// GetType 获取类型
func (x *Source) GetType() int64 {
	return ConnTypeSource
}

// Serialize 序列化成[]byte
func (x *Source) Serialize() []byte {
	b, err := msgpack.Marshal(x)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// Deserialize 反序列化
func (x *Source) Deserialize(d []byte) error {
	err := msgpack.Unmarshal(d, x) // 将二进制流转化回结构体
	if err != nil {
		return err
	}
	return nil
}

// Check 校验
func (x *Source) Check(cont *blc.Cont) bool {
	if x.From != blc.GetAddressFromPublicKey(cont.PublicKey) {
		return false
	}
	if x.Content == "" {
		return false
	}
	return true
}
