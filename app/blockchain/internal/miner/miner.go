package miner

import (
	"github.com/google/wire"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/common"
	"bigcat_test_coin/tools"
)

//最大校验块儿可靠性的人数
const maxVerifierMinerNumber = 10

var ProviderSet = wire.NewSet(NewMinerControl)

type Miner struct {
	//地址
	Addr string
	//编号
	Number int
	//高度
	LastHeight  int64
	LoseHeight  int64
	ClearHeight int64
}

type MinerControl struct {
	Miners []Miner
}

func NewMinerControl() *MinerControl {
	return &MinerControl{Miners: make([]Miner, 0, 8)}
}

func (m *MinerControl) AddMiner(height int64, addr string) {
	for _, miner := range m.Miners {
		if miner.Addr == addr {
			return
		}
	}
	m.Miners = append(m.Miners, Miner{
		Addr:       addr,
		Number:     len(m.Miners) + 1,
		LastHeight: height,
	})
}
func (m *MinerControl) Lose(height int64, miners ...string) {
	clearHeight := height + height%int64(common.BlockChainConsensusMechanismSize) + int64(common.BlockChainConsensusMechanismSize)
	for _, addr := range miners {
		for index, miner := range m.Miners {
			if miner.Addr == addr {
				if con := m.Miners[index].LoseHeight; con == 0 || con > height {
					m.Miners[index].LoseHeight = height
					m.Miners[index].ClearHeight = clearHeight
				}
				break
			}
		}
	}
}
func (m *MinerControl) Recover(height int64) {
	ms := make([]Miner, 0, len(m.Miners))
	index := 1
	for _, miner := range m.Miners {
		if miner.LastHeight > height {
			continue
		}
		if miner.LoseHeight >= height {
			miner.LoseHeight = 0
			miner.ClearHeight = 0
		}
		miner.Number = index
		index++
	}
	m.Miners = ms
}
func (m *MinerControl) Clear(height int64) {
	if height%int64(common.BlockChainConsensusMechanismSize) != 0 {
		return
	}
	var index int = 1
	ms := make([]Miner, 0, len(m.Miners))
	for _, miner := range m.Miners {
		if miner.ClearHeight < height {
			miner.Number = index
			index++
		}
	}
	m.Miners = ms
}
func (m *MinerControl) Truncation(height int64) *MinerControl {
	newControl := NewMinerControl()
	for _, miner := range m.Miners {
		if miner.LastHeight < height {
			newControl.Miners = append(newControl.Miners, miner)
		} else {
			break
		}
	}
	return newControl
}
func (m *MinerControl) GetVerifierMiners(hx []byte) (miners []string) {
	if len(m.Miners) == 0 {
		return
	}
	miners = make([]string, 0)
	divisor := len(m.Miners)
	divisor = tools.NumberSqrt(divisor)
	if divisor > maxVerifierMinerNumber {
		divisor = maxVerifierMinerNumber
	}
	_, numberRem := tools.NumberQuoAndRem(hx, divisor)
	for _, v := range m.Miners {
		if numberRem == 0 {
			miners = append(miners, v.Addr)
			continue
		}
		if v.Number%numberRem == 0 {
			miners = append(miners, v.Addr)
		}
	}
	return miners
}
func (m *MinerControl) GetMinerCount() int {
	return len(m.Miners)
}

func (m *MinerControl) GetMiner(addr string) Miner {
	for _, miner := range m.Miners {
		if miner.Addr == addr {
			return miner
		}
	}
	return Miner{}
}
func (m *MinerControl) CheckMiner(addr string) bool {
	if m.GetMinerCount() == 0 {
		return true
	}
	miner := m.GetMiner(addr)
	if miner.Number > 0 {
		return true
	}
	return false
}

// GetNewMiners 获取新增的矿工地址
func (m *MinerControl) GetNewMiners(block *blc.Block) (miners []string) {
	miners = make([]string, 0)
	for _, t := range block.Transactions {
		if t.To == common.BlockChainMinerAddr {
			miners = append(miners, t.GetForm())
		}
	}
	return miners
}

// CheckMinerBlock 校验矿工协议
func (m *MinerControl) CheckMinerBlock(block *blc.Block) bool {
	for _, tran := range block.Transactions {
		if tran.To == common.BlockChainMinerAddr {
			if tran.Value != common.BlockChainMinerFee {
				return false
			}
		}
	}
	return true
}
