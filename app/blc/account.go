package blc

import (
	"google.golang.org/protobuf/proto"
)



// Serialize 将account序列化成[]byte
func (x *Account) Serialize() []byte {
	bys, err := proto.Marshal(x)
	if err != nil {
		panic(err.Error())
	}
	return bys
}

func (x *Account) Deserialize(d []byte) {
	err := proto.Unmarshal(d, x)
	if err != nil {
		panic(err.Error())
	}
}
