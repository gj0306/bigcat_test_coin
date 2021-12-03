package chain

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	v1 "bigcat_test_coin/app/api"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/common"
	"bigcat_test_coin/tools"
	"time"
)


func N2PGetInfo(c *Chain) {
	fun := func(nid, data []byte) error {
		req := &v1.N2PGetInfo{}
		err := proto.Unmarshal(data, req)
		if err != nil {
			return err
		}
		leftBlock := c.cBox.GetBlock(c.cBox.LeftHeight)
		LastBlock := c.cBox.GetBlock(c.cBox.RightHeight)
		if req.SearchHeight == 0{
			req.SearchHeight = 1
		}
		searchBlock := c.cBox.GetBlock(req.SearchHeight)
		reply := &v1.P2NGetInfo{
			Height:   c.cBox.RightHeight,
			LeftHash: leftBlock.Hash,
			LastHash: LastBlock.Hash,
			SearchHash: searchBlock.Hash,
			Now:      time.Now().UnixMilli(),
			VerifierTotal: LastBlock.VerifierTotal,
		}
		bys, err := proto.Marshal(reply)
		if err != nil {
			return err
		}
		err = c.work.SendOneMessage(v1.EncodingMessage(v1.P2NGetInfoNumber, bys), nid)
		if err != nil {
			return err
		}
		return nil
	}
	c.recvApp[v1.N2PGetInfoNumber] = fun
}

func P2NGetInfo(c *Chain) {
	fun := func(nid, data []byte) error {
		req := &v1.P2NGetInfo{}
		err := proto.Unmarshal(data, req)
		if err != nil {
			return err
		}
		if c.nodeInfos != nil {
			c.nodeInfos[string(nid)] = req
		}
		return nil
	}
	c.recvApp[v1.P2NGetInfoNumber] = fun
}

func N2PGetBlock(c *Chain){
	fun := func(nid, data []byte) error {
		req := &v1.N2PGetBlock{}
		err := proto.Unmarshal(data, req)
		if err != nil {
			return err
		}
		lashBlock := c.cBox.GetBlockByHash(req.GetLastHash())
		var bys []byte
		var num uint8
		var count int64 = 10
		var left int64 = req.GetRightHeight()+1
		switch lashBlock.Height {
		case 0:
			if req.GetRightHeight() >= int64(common.BlockChainConsensusMechanismSize){
				//不在范围内
				left = req.GetLeftHeight()
				right := req.GetRightHeight()
				num = v1.N2PGetBlocksHashNumber
				reply := &v1.N2PGetBlocksHash{
					BlockTxs: make([][]byte,0),
					LeftHeight: left,
				}
				for left<=right{
					block := c.cBox.GetBlock(left)
					if block.Height == 0{
						break
					}
					reply.BlockTxs = append(reply.BlockTxs, block.Hash)
					left++
				}
				bys, err = proto.Marshal(reply)
				if err != nil {
					return err
				}
				break
			}
			left = 1
			count = 60
			fallthrough
		default:
			num = v1.P2NGetBlockNumber
			reply := &v1.P2NGetBlock{
				Blocks: make([]*blc.Block,0),
			}
			right := left + count
			for left<=right{
				block := c.cBox.GetBlock(left)
				if block == nil{
					break
				}
				if block.Height == 0{
					break
				}
				reply.Blocks = append(reply.Blocks, block)
				left++
			}
			bys, err = proto.Marshal(reply)
			if err != nil {
				return err
			}
		}
		err = c.work.SendOneMessage(v1.EncodingMessage(num, bys), nid)
		if err != nil {
			return err
		}
		return nil
	}
	c.recvApp[v1.N2PGetBlockNumber] = fun
}

func N2PGetBlocksHash(c *Chain){
	fun := func(nid, data []byte) error {
		req := &v1.N2PGetBlocksHash{}
		err := proto.Unmarshal(data, req)
		if err != nil {
			return err
		}
		leftBlock := c.cBox.GetBlock(c.cBox.LeftHeight)
		reply := &v1.N2PGetBlock{
			LeftHeight: leftBlock.Height,
			LastHash: leftBlock.Hash,
		}
		for index,hx := range req.GetBlockTxs(){
			height := int64(index) + req.LeftHeight
			block := c.cBox.GetBlock(height)
			if !tools.BytesEqual(hx,block.Hash){
				reply.RightHeight = block.Height -1
				break
			}
		}
		bys, err := proto.Marshal(reply)
		if err != nil {
			return err
		}
		err = c.work.SendOneMessage(v1.EncodingMessage(v1.N2PGetBlockNumber, bys), nid)
		return nil
	}
	c.recvApp[v1.N2PGetBlocksHashNumber] = fun
}

func P2NGetBlock(c *Chain){
	fun := func(nid, data []byte) error {
		if !tools.BytesEqual(nid,c.nodeId){
			return fmt.Errorf("不认可的nid数据")
		}
		defer func() {
			c.nodeId = nil
		}()
		req := &v1.P2NGetBlock{}
		err := proto.Unmarshal(data, req)
		if err != nil {
			return err
		}
		err = c.cBox.AddBlocks(req.Blocks)
		return err
	}
	c.recvApp[v1.P2NGetBlockNumber] = fun
}

func PlushVerifier(c *Chain){
	fun := func(nid, data []byte) error {
		req := &v1.PlushVerifier{}
		err := proto.Unmarshal(data, req)
		if err != nil {
			return err
		}
		if req.Verifier == nil{
			return nil
		}
		c.cBox.AddVerifier(req.Verifier)
		return nil
	}
	c.recvApp[v1.PlushVerifierNumber] = fun
}

func PushContracts(c *Chain)  {
	fun := func(nid, data []byte) error {
		req := &v1.PushContracts{}
		err := proto.Unmarshal(data, req)
		if err != nil {
			return err
		}
		if req.Conts == nil{
			return nil
		}
		for _,cont := range req.GetConts(){
			c.cBox.AddContract(cont)
		}
		return nil
	}
	c.recvApp[v1.PushContractsNumber] = fun
}

func PushTransaction(c *Chain)  {
	fun := func(nid, data []byte) error {
		req := &v1.PushTransactions{}
		err := proto.Unmarshal(data, req)
		if err != nil {
			return err
		}
		if req.Trans == nil{
			return nil
		}
		for _,tran := range req.GetTrans(){
			c.cBox.AddTransaction(tran)
		}
		return nil
	}
	c.recvApp[v1.PushTransactionsNumber] = fun
}