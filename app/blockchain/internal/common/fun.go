package common

import (
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/tools"
	"time"
)

// GetNextBlockTime 获取下一发块儿时间
func GetNextBlockTime(lastTm int64, delay int) int64 {
	lastTime := time.Unix(lastTm, 0)
	tm := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), lastTime.Hour(), lastTime.Minute(), 0, 0, time.Local)
	return tm.Add(time.Minute * time.Duration(delay)).Unix()
}

// BlockWeightRatio 块儿优先级
func BlockWeightRatio(b1,b2 *blc.Block)bool{
	if b1.VerifierTotal > b2.VerifierTotal{
		return true
	}else if b1.VerifierTotal < b2.VerifierTotal {
		return false
	}
	//if len(b1.Transactions)>len(b2.Transactions){
	//	return true
	//}else if len(b1.Transactions)<len(b2.Transactions) {
	//	return false
	//}
	return tools.BytesCmp(b1.Hash,b2.Hash)
}

// GetContractFee 合约价格计算
func GetContractFee(cont *blc.Cont)int64{
	num := len(cont.Data)
	return int64(num/ContractFeeBytes) + ContractFee
}