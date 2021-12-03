package db

import (
	blc "bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/config"
	"bigcat_test_coin/app/blockchain/internal/db/badger_control"
	"bigcat_test_coin/app/blockchain/internal/db/boltdb_control"
	"bigcat_test_coin/tools"
	"fmt"
	"github.com/google/wire"
	"go.uber.org/zap"
	"os"
)

type Control interface {
	LoadConfig(config config.DbConf) error
	GetBlock(height int64) *blc.Block
	GetBlocks(leftHeight, rightHeight int64) (blocks []*blc.Block)
	GetBlockByHash(hash []byte) (block *blc.Block)
	GetContract(tx []byte, height int64) *blc.Cont
	GetTransaction(tx []byte, height int64) *blc.Transaction
	SaveBlock(block *blc.Block)
	GetHeight() int64
	GetAccount(addr string) *blc.Account
	SaveAccount(addr string, account *blc.Account)
	GetTransactionHeight(tx []byte) int64
	GetContractHeight(tx []byte) int64
}

var ProviderSet = wire.NewSet(NewDbControl)

func NewDbControl(conf *config.DbConf, log *zap.Logger) (dbControl Control, err error) {
	if conf.Dir == "" {
		conf.Dir = "./db"
	}
	ok, _ := tools.IsFileExist(conf.Dir)
	if !ok {
		err = os.Mkdir(conf.Dir, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}
	switch conf.DriverName {
	case "boltdb":
		dbControl = boltdb_control.NewDb(conf.Dir + "/")
	case "badger":
		dbControl = badger_control.NewDb(conf.Dir + "/")
	default:
		return nil, fmt.Errorf("unknown DB DriverName " + conf.DriverName)
	}
	err = dbControl.LoadConfig(*conf)
	if err != nil {
		return nil, err
	}
	return dbControl, nil
}
