package common

// BlockChainConsensusMechanismSize 最大合约认证数
const BlockChainConsensusMechanismSize int = 60

// BlockChainMinerAddr 矿工地址
const BlockChainMinerAddr = "1M8QQfDtxJvNsE3Cw7SXA551zndDzwJT2T"

// BlockChainMinerFee 矿工价格
const BlockChainMinerFee int64 = 10


const (
	//WorkIncome 工作收益
	WorkIncome int64 = 10
	// TransactionFee
	//交易手续费 fee = value / TransactionFee
	//如果费用小于 TransactionFee 交易失败
	TransactionFee int64 = 10

	// ContractFee 合约价格
	ContractFee int64 = 1
	ContractFeeBytes = 512
)