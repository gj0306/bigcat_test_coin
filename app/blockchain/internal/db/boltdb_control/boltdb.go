package boltdb_control

import /**/ (
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/config"
	"bigcat_test_coin/tools"
)

// BoltDbControl 数据库
type BoltDbControl struct {
	fileDb                  *BlockchainDB
	info                    *basicInfo
	dBFileName              string
}
func NewDb(path string) *BoltDbControl {
	db := &BoltDbControl{
		fileDb:                  New(),
		info:                    &basicInfo{},
		dBFileName:              path + "data.db",
	}
	db.fileDb.DeleteBucket(db.dBFileName, AccountBucket)
	//db.fileDb.DeleteBucket(db.dBFileName,ContractsInfoBucket)
	return db
}

func (db BoltDbControl) LoadConfig(config config.DbConf)error{
	info := db.getInfo()
	db.info.Height = info.Height
	return nil
}

func (db BoltDbControl) GetBlock(height int64) *blc.Block {
	key := tools.Int64ToBytes(height)
	data := db.fileDb.View(db.dBFileName,key, BlocksInfoBucket)
	if len(data)>0{
		info := blockInfo{}
		info.Deserialize(data)
		return db.GetBlockByHash(info.HxHash)
	}
	return nil
}
func (db BoltDbControl) GetBlocks(leftHeight, rightHeight int64) (blocks []*blc.Block) {
	blocks = make([]*blc.Block,0)
	for leftHeight<=rightHeight{
		block := db.GetBlock(leftHeight)
		if block == nil{
			break
		}
		if block.Height == 0{
			break
		}
		leftHeight++
		blocks = append(blocks, block)
	}
	return blocks
}
func (db BoltDbControl) GetBlockByHash(hash []byte)(block *blc.Block) {
	data := db.fileDb.View(db.dBFileName,hash, BlockBucket)
	if len(data)>0{
		block = &blc.Block{}
		block.Deserialize(data)
		return
	}
	return
}
func (db BoltDbControl) GetContract(tx []byte,height int64) *blc.Cont {
	data := db.fileDb.View(db.dBFileName,tx, ContractsInfoBucket)
	if len(data)>0{
		cInfo := contractInfo{}
		cInfo.Deserialize(data)
		block := db.GetBlockByHash(cInfo.HxHash)
		if block.Height>0&&block.Height<=height{
			return block.Contracts[cInfo.Index]
		}
	}
	return nil
}
func (db BoltDbControl) GetTransaction(tx []byte,height int64) *blc.Transaction {
	data := db.fileDb.View(db.dBFileName,tx, TransactionInfoBucket)
	if len(data)>0{
		cInfo := transactionInfo{}
		cInfo.Deserialize(data)
		block := db.GetBlockByHash(cInfo.HxHash)
		if block.Height>0&&block.Height<=height{
			return block.Transactions[cInfo.Index]
		}
	}
	return nil
}
func (db BoltDbControl) GetTransactionHeight(tx []byte)int64{
	data := db.fileDb.View(db.dBFileName,tx, TransactionInfoBucket)
	if len(data)>0{
		cInfo := transactionInfo{}
		cInfo.Deserialize(data)
		return cInfo.Height
	}
	return 0
}
func (db BoltDbControl) GetContractHeight(tx []byte)int64{
	data := db.fileDb.View(db.dBFileName,tx, ContractsInfoBucket)
	if len(data)>0{
		cInfo := contractInfo{}
		cInfo.Deserialize(data)
		return cInfo.Height
	}
	return 0
}
func (db BoltDbControl) SaveBlock(block *blc.Block)  {
	if block.Height > db.info.Height+1 {
		return
	}
	switch block.Height {
	default:
		db.truncationBlockChain(block.Height)
		fallthrough
	case  db.info.Height+1:
		db.saveBlock(block)
	}
	db.info.Height = block.Height
	db.saveInfo()
	return
}


func (db BoltDbControl) truncationBlockChain(height int64){
	for height<=db.info.Height{
		db.delBlockByHeight(height)
		height++
	}
}
func (db BoltDbControl) saveBlock(block *blc.Block)  {
	//保存索引
	info := blockInfo{
		Height: block.Height,
		HxHash: block.Hash,
	}
	key := tools.Int64ToBytes(block.Height)
	db.fileDb.Put(db.dBFileName,key,info.Serialize(), BlocksInfoBucket)
	//保存块儿
	db.fileDb.Put(db.dBFileName,block.Hash,block.Serialize(), BlockBucket)
	//保存合约索引
	for index, conn := range block.Contracts{
		cInfo := contractInfo{
			TxHash: conn.TxHash,
			HxHash: block.Hash,
			Index:  index,
			Height: block.Height,
		}
		db.fileDb.Put(db.dBFileName, conn.TxHash,cInfo.Serialize(), ContractsInfoBucket)
	}
	//保存交易索引
	for index, transaction := range block.Transactions{
		cInfo := transactionInfo{
			TxHash: transaction.TxHash,
			HxHash: block.Hash,
			Height: block.Height,
			Index:  index,
		}
		db.fileDb.Put(db.dBFileName, transaction.TxHash,cInfo.Serialize(), TransactionInfoBucket)
	}
}
func (db BoltDbControl) delBlockByHeight(height int64)  {
	block := db.GetBlock(height)
	if block.Height==0{
		return
	}
	db.fileDb.Delete(db.dBFileName,block.Hash, BlockBucket)
	for _, conn := range block.Contracts{
		db.fileDb.Delete(db.dBFileName, conn.TxHash, ContractsInfoBucket)
	}
	for _, t := range block.Transactions{
		db.fileDb.Delete(db.dBFileName, t.TxHash, TransactionInfoBucket)
	}
}

func (db BoltDbControl) GetHeight()int64 {
	return db.info.Height
}

func (db BoltDbControl) GetAccount(addr string) (block *blc.Account) {
	data := db.fileDb.View(db.dBFileName,[]byte(addr), AccountBucket)
	if len(data)>0{
		account := &blc.Account{}
		account.Deserialize(data)
		return account
	}
	return nil
}
func (db BoltDbControl) SaveAccount(addr string,account *blc.Account) {
	//bys := make([]byte,8)
	//binary.BigEndian.PutUint64(bys,uint64(account))
	db.fileDb.Put(db.dBFileName,[]byte(addr),account.Serialize(), AccountBucket)
}

//info
func (db BoltDbControl) getInfo() (basic basicInfo){
	data := db.fileDb.View(db.dBFileName,[]byte("info"), BasicBucket)
	if len(data)>0{
		basic.Deserialize(data)
		return
	}
	return
}
func (db BoltDbControl) saveInfo(){
	db.fileDb.Put(db.dBFileName,[]byte("info"),db.info.Serialize(), BasicBucket)
}
