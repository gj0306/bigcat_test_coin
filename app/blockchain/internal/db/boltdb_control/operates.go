/*
	本包是作为对blot数据库封装的一个存在
*/
package boltdb_control

import (
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"sync"
)

var mu sync.Mutex

var dbMap map[string]*bolt.DB = map[string]*bolt.DB{}

func getDb(dBFileName string) *bolt.DB{
	mu.Lock()
	defer func() {
		mu.Unlock()
	}()
	if db,ok := dbMap[dBFileName];ok{
		return db
	}
	newDb, err := bolt.Open(dBFileName, 0600, nil)
	if err != nil{
		log.Panic(err.Error())
	}
	dbMap[dBFileName] = newDb
	return newDb
}

// Put 存入数据
func (bd *BlockchainDB) Put(dBFileName string,k, v []byte, bt BucketType) {
	db := getDb(dBFileName)
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			var err error
			bucket, err = tx.CreateBucket([]byte(bt))
			if err != nil {
				log.Panic(err.Error())
			}
		}
		err := bucket.Put(k, v)
		if err != nil {
			log.Panic(err.Error())
		}
		return nil
	})
	if err != nil {
		log.Panic(err.Error())
	}
}

// View 查看数据
func (bd *BlockchainDB) View(dBFileName string,k []byte, bt BucketType) []byte {
	var err error
	db := getDb(dBFileName)
	result := []byte{}
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			msg := "datebase view warnning:没有找到仓库：" + string(bt)
			return errors.New(msg)
		}
		result = bucket.Get(k)
		return nil
	})
	if err != nil {
		//log.Warn(err)
		return nil
	}
	//不再次赋值的话，返回值会报错，不知道狗日的啥意思
	realResult := make([]byte, len(result))
	copy(realResult, result)
	return realResult
}
// Delete 删除数据
func (bd *BlockchainDB) Delete(dBFileName string,k []byte, bt BucketType) bool {
	var err error
	db := getDb(dBFileName)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			msg := "datebase delete warnning:没有找到仓库：" + string(bt)
			return errors.New(msg)
		}
		err = bucket.Delete(k)
		if err != nil {
			log.Panic(err.Error())
		}
		return nil
	})
	if err != nil {
		log.Panic(err.Error())
	}
	return true
}

// DeleteBucket 删除仓库
func (bd *BlockchainDB) DeleteBucket(dBFileName string,bt BucketType) bool {
	var err error
	db := getDb(dBFileName)
	err = db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bt))
	})
	if err != nil {
		switch err.Error() {
		case "bucket not found":
		default:
			log.Panic(err.Error())
		}
	}
	return true
}
