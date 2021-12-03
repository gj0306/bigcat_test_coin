// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/box"
	"bigcat_test_coin/app/blockchain/internal/chain"
	"bigcat_test_coin/app/blockchain/internal/coinlog"
	"bigcat_test_coin/app/blockchain/internal/config"
	"bigcat_test_coin/app/blockchain/internal/db"
	"bigcat_test_coin/app/blockchain/internal/miner"
	"bigcat_test_coin/app/blockchain/internal/network"
	"bigcat_test_coin/app/blockchain/internal/server"
	"bigcat_test_coin/app/blockchain/internal/service"
)

// initApp init kratos application.
func initApp(privateKey *blc.Address,conf *config.YamlConfig,logConf *config.LogConf,dbConf *config.DbConf,addressConf *config.AddressConf,networkConf *config.NetworkConf,serviceConf *config.WebServer,log log.Logger) (*chain.Chain, func(), error) {
	panic(wire.Build(
		box.ProviderSet,
		db.ProviderSet,
		network.ProviderSet,
		miner.ProviderSet,
		coinlog.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
		newApp))
}
