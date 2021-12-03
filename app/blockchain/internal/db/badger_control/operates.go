/*
	本包是作为对blot数据库封装的一个存在
*/
package badger_control

import (
	badger "github.com/dgraph-io/badger/v3"
	"log"
)

// Put 存入数据
func Put(db *badger.DB,k, v []byte) {
	var err error
	err = db.Update(func(txn *badger.Txn) error {
		err = txn.Set(k,v)
		return err
	})
	if err != nil{
		log.Println(err.Error())
	}
}

// View 查看数据
func View(db *badger.DB,k []byte) []byte {
	var err error
	var data []byte
	err = db.View(func(txn *badger.Txn) error {
		item,err := txn.Get(k)
		if err == nil && item != nil{
			err = item.Value(func(val []byte) error {
				data = val
				return nil
			})
		}
		return err
	})
	if err != nil{
		log.Println(string(k),err.Error())
	}
	return data
}

func Views(db *badger.DB,ks [][]byte) [][]byte {
	var err error
	data := make([][]byte,len(ks))
	err = db.View(func(txn *badger.Txn) error {
		for i,k := range ks{
			item,err := txn.Get(k)
			if err == nil && item != nil{
				err = item.Value(func(val []byte) error {
					data[i] = val
					return nil
				})
			}
			return err
		}
		return nil
	})
	if err != nil{
		log.Println(err.Error())
	}
	return data
}

// Delete 删除数据
func Delete(db *badger.DB,k []byte) bool {
	var err error
	err = db.Update(func(txn *badger.Txn) error {
		err = txn.Delete(k)
		return err
	})
	if err != nil{
		log.Println(err.Error())
		return false
	}
	return true
}

