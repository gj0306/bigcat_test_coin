package contract

import (
	"github.com/vmihailenco/msgpack"
	"bigcat_test_coin/app/blc"
)

func CreateNews(title,content string) News {
	return News{
		Title:   title,
		Content: content,
	}
}

// GetType 获取类型
func (x *News) GetType() int64 {
	return ConnTypeNews
}

// Serialize 序列化成[]byte
func (x *News) Serialize() []byte {
	b, err := msgpack.Marshal(x)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// Deserialize 反序列化
func (x *News) Deserialize(d []byte) error {
	err := msgpack.Unmarshal(d, x) // 将二进制流转化回结构体
	if err != nil {
		return err
	}
	return nil
}

// Check 校验
func (x *News) Check(cont *blc.Cont) bool {
	if len(x.Title)==0 || len(x.Content)==0{
		return false
	}
	return true
}
