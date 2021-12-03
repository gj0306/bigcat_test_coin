/*
	本包是作为对blot数据库封装的一个存在
*/
package boltdb_control

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
)

var ListenPort string

// BucketType 仓库类型
type BucketType string

const (
	BlockBucket BucketType = "blocks"
	BasicBucket BucketType = "basic"

	/*临时表*/
	AccountBucket         BucketType = "account"
	ContractsInfoBucket   BucketType = "contracts"
	TransactionInfoBucket BucketType = "transactions"
	BlocksInfoBucket      BucketType = "block_info"
)

type BlockchainDB struct {
	ListenPort string
}

func New() *BlockchainDB {
	bd := &BlockchainDB{ListenPort}
	return bd
}

// IsBlotExist 判断数据库是否存在
func IsBlotExist(nodeID string) bool {
	var DBFileName = "news_" + nodeID + ".coindb"
	_, err := os.Stat(DBFileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// IsBucketExist 判断仓库是否存在
func IsBucketExist(bd *BlockchainDB, bt BucketType) bool {
	var isBucketExist bool
	var DBFileName = "news_" + ListenPort + ".coindb"
	db, err := bolt.Open(DBFileName, 0600, nil)
	if err != nil {
		log.Panic(err.Error())
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			isBucketExist = false
		} else {
			isBucketExist = true
		}
		return nil
	})
	if err != nil {
		log.Panic("datebase IsBucketExist err:" + err.Error())
	}

	err = db.Close()
	if err != nil {
		log.Panic("coindb close err :" + err.Error())
	}
	return isBucketExist
}
