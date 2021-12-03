package boltdb_control

import (
	"github.com/vmihailenco/msgpack"
)

type basicInfo struct {
	Height      int64
}

func (b *basicInfo) Serialize() []byte {
	data, err := msgpack.Marshal(b)
	if err != nil {
		panic(err.Error())
	}
	return data
}
func (b *basicInfo) Deserialize(d []byte) {
	err := msgpack.Unmarshal(d, b)
	if err != nil {
		panic(err.Error())
		return
	}
}

type blockInfo struct {
	Height int64
	HxHash []byte
}

func (b *blockInfo) Serialize() []byte {
	data, err := msgpack.Marshal(b)
	if err != nil {
		panic(err.Error())
	}
	return data
}
func (b *blockInfo) Deserialize(d []byte) {
	err := msgpack.Unmarshal(d, b)
	if err != nil {
		panic(err.Error())
		return
	}
}

type contractInfo struct {
	TxHash []byte
	HxHash []byte
	Height int64
	Index  int
}

func (c *contractInfo) Serialize() []byte {
	data, err := msgpack.Marshal(c)
	if err != nil {
		panic(err.Error())
	}
	return data
}
func (c *contractInfo) Deserialize(d []byte) {
	err := msgpack.Unmarshal(d, c)
	if err != nil {
		panic(err.Error())
		return
	}
}

type transactionInfo struct {
	TxHash []byte
	HxHash []byte
	Height int64
	Index  int
}
func (t *transactionInfo) Serialize() []byte {
	data, err := msgpack.Marshal(t)
	if err != nil {
		panic(err.Error())
	}
	return data
}
func (t *transactionInfo) Deserialize(d []byte) {
	err := msgpack.Unmarshal(d, t)
	if err != nil {
		panic(err.Error())
		return
	}
}