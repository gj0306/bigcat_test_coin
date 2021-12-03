package badger_control

import (
	badger "github.com/dgraph-io/badger/v3"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/config"
	"bigcat_test_coin/tools"
	"os"
	"strconv"
	"time"
)

/*
注意事项徽章不直接使用Cgo，但它依赖于https://github.com/DataDog/zstd进行压缩，它需要 gcc / cgo 。
如果您希望使用没有 gcc/cgo 的徽章，则可以运行该徽章，无需支持 ZSTD 压缩算法即可下载徽章,那么请设置 CGO_ENABLED=0。

*/

const debug = false
const (
	BlockBucket = "blocks"
	BasicBucket = "basic"
	// TableIndex 索引表
	TableIndex = "index"
	// AccountBucket 临时表
	AccountBucket = "account"
)

const (
	BlockIndex uint8 = 1
	TranIndex  uint8 = 2
	ContIndex  uint8 = 3
)

// BadgerControl 数据库
type BadgerControl struct {
	dbs                     map[string]*badger.DB
	height                  int64
	path string
}

func NewDb(path string) *BadgerControl {
	bdb := &BadgerControl{
		dbs:                     map[string]*badger.DB{},
		height:                  time.Now().Unix(),
		path:                    path,
	}
	bdb.dbs[BlockBucket] = newDb(path,BlockBucket)
	bdb.dbs[BasicBucket] = newDb(path,BasicBucket)
	bdb.dbs[TableIndex] = newDb(path,TableIndex)

	delDb(path,AccountBucket)
	bdb.dbs[AccountBucket] = newDb(path,AccountBucket)

	height := bdb.getHeight()
	bdb.height = height
	return bdb
}

func newDb(path,name string) *badger.DB {
	var opt badger.Options
	if debug {
		opt = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opt = badger.DefaultOptions(path + name)
	}
	db, err := badger.Open(opt)
	if err != nil {
		panic(err.Error())
	}
	return db
}
func delDb(path,name string) {
	_,err := os.Stat(path+name)
	if err != nil{
		if os.IsNotExist(err){
			return
		}
		panic(err.Error())
	}
	err = os.RemoveAll(path + name)
	if err != nil {
		panic(err.Error())
	}
}

func (c *BadgerControl) LoadConfig(config config.DbConf) error {
	return nil
}
func (c *BadgerControl) GetBlock(height int64) *blc.Block {
	return c.getBlock(height)
}
func (c *BadgerControl) GetBlocks(leftHeight, rightHeight int64) (blocks []*blc.Block) {
	blocks = make([]*blc.Block, 0)
	for leftHeight <= rightHeight {
		block := c.getBlock(leftHeight)
		if block == nil {
			break
		}
		leftHeight++
		blocks = append(blocks, block)
	}
	return blocks
}
func (c *BadgerControl) GetBlockByHash(hash []byte) (block *blc.Block) {
	return c.getBlockByHash(hash)
}

func (c *BadgerControl) GetContract(tx []byte, height int64) *blc.Cont {
	h := c.GetContractHeight(tx)
	if h <= 0 || h>height{
		return nil
	}
	block := c.getBlock(h)
	if block == nil {
		return nil
	}
	for _, cont := range block.Contracts {
		if tools.BytesEqual(cont.TxHash, tx) {
			return cont
		}
	}
	return nil
}
func (c *BadgerControl) GetTransaction(tx []byte, height int64) *blc.Transaction {
	h := c.GetTransactionHeight(tx)
	if h <= 0 || h>height{
		return nil
	}
	block := c.getBlock(h)
	if block == nil {
		return nil
	}
	for _, tran := range block.Transactions {
		if tools.BytesEqual(tran.TxHash, tx) {
			return tran
		}
	}
	return nil
}
func (c *BadgerControl) SaveBlock(block *blc.Block) {
	switch block.Height {
	default:
		c.truncationBlockChain(block.Height)
		fallthrough
	case c.height + 1:
		c.saveBlock(block)
	}
	c.height = block.Height
	c.saveHeight()
	return
}
func (c *BadgerControl) truncationBlockChain(height int64) {
	for height <= c.height {
		c.delBlock(height)
		height++
	}
}
func (c *BadgerControl) GetHeight() int64 {
	return c.height
}
func (c *BadgerControl) GetAccount(addr string) (block *blc.Account) {
	val := View(c.dbs[AccountBucket], []byte(addr))
	if len(val) == 0 {
		return nil
	}
	acc := &blc.Account{}
	acc.Deserialize(val)
	return acc
}
func (c *BadgerControl) SaveAccount(addr string, account *blc.Account) {
	Put(c.dbs[AccountBucket], []byte(addr), []byte(account.Serialize()))
}

func (c *BadgerControl) GetTransactionHeight(tx []byte)int64{
	return c.getIndex(TranIndex, tx)
}
func (c *BadgerControl) GetContractHeight(tx []byte)int64{
	return c.getIndex(ContIndex, tx)
}

func (c *BadgerControl) getHeight() int64 {
	val := View(c.dbs[BasicBucket], []byte{0})
	height, _ := strconv.ParseInt(string(val), 36, 64)
	return height
}
func (c *BadgerControl) saveHeight() {
	val := strconv.FormatInt(c.height, 36)
	Put(c.dbs[BasicBucket], []byte{0}, []byte(val))
}
func (c *BadgerControl) getBlock(height int64) *blc.Block {
	key := strconv.FormatInt(height, 36)
	val := View(c.dbs[BlockBucket], []byte(key))
	if len(val) == 0 {
		return nil
	}
	block := &blc.Block{}
	block.Deserialize(val)
	return block
}
func (c *BadgerControl) getBlockByHash(hash []byte) *blc.Block {
	val := View(c.dbs[TableIndex], append([]byte{BlockIndex}, hash...))
	if len(val) == 0 {
		return nil
	}
	height, _ := strconv.ParseInt(string(val), 36, 64)
	return c.getBlock(height)
}
func (c *BadgerControl) saveBlock(block *blc.Block) {
	key := []byte(strconv.FormatInt(block.Height, 36))
	Put(c.dbs[BlockBucket], key, block.Serialize())
	//索引相关保存
	Put(c.dbs[TableIndex], append([]byte{BlockIndex}, block.Hash...), block.Serialize())
	//交易
	for _, tran := range block.Transactions {
		Put(c.dbs[TableIndex], append([]byte{TranIndex}, tran.TxHash...), key)
	}
	//合约
	for _, cont := range block.Contracts {
		Put(c.dbs[TableIndex], append([]byte{ContIndex}, cont.TxHash...), key)
	}
}
func (c *BadgerControl) delBlock(height int64) {
	block := c.getBlock(height)
	if block == nil {
		return
	}
	key := strconv.FormatInt(height, 36)
	Delete(c.dbs[BlockBucket], []byte(key))
	//删除索引
	Delete(c.dbs[TableIndex], append([]byte{BlockIndex}, block.Hash...))
	//删除交易
	for _, tran := range block.Transactions {
		Delete(c.dbs[TableIndex], append([]byte{TranIndex}, tran.TxHash...))
	}
	//删除合约
	for _, cont := range block.Contracts {
		Delete(c.dbs[TableIndex], append([]byte{ContIndex}, cont.TxHash...))
	}
}
func (c *BadgerControl) getIndex(index uint8, key []byte) int64 {
	val := View(c.dbs[TableIndex], append([]byte{index}, key...))
	if len(val) == 0 {
		return 0
	}
	num, _ := strconv.ParseInt(string(val), 36, 64)
	return num
}
