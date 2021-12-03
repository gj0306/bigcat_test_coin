package v1

const (
	N2PGetInfoNumber       uint8 = 1
	P2NGetInfoNumber       uint8 = 2
	N2PGetBlockNumber      uint8 = 3
	P2NGetBlockNumber      uint8 = 4
	N2PGetBlocksHashNumber uint8 = 5

	PushBlockNumber     uint8 = 101
	PlushVerifierNumber uint8 = 102
	PushContractsNumber uint8 = 103
	PushTransactionsNumber uint8 = 104
)

// EncodingMessage 编码消息
func EncodingMessage(number uint8,data []byte) []byte {
	bys := make([]byte,0,len(data)+1)
	bys = append(bys, number)
	bys = append(bys, data...)
	return bys
}

// DecodingMessage 解码消息
func DecodingMessage(data []byte) (number uint8,vals []byte) {
	if len(data) == 0{
		return 0, nil
	}
	return data[0], data[1:]
}