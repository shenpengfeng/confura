package cmd

import (
	sdk "github.com/Conflux-Chain/go-conflux-sdk"
	"github.com/conflux-chain/conflux-infura/store"
	"github.com/conflux-chain/conflux-infura/store/mysql"
	"github.com/conflux-chain/conflux-infura/store/redis"
	"github.com/conflux-chain/conflux-infura/util/rpc"
	"github.com/openweb3/web3go"
)

type storeContext struct {
	cfxDB    *mysql.MysqlStoreV2
	ethDB    *mysql.MysqlStoreV2
	cfxCache *redis.RedisStore
}

func mustInitStoreContext() storeContext {
	var ctx storeContext

	if config := mysql.MustNewConfigFromViper(); config.Enabled {
		ctx.cfxDB = config.MustOpenOrCreate(mysql.StoreOption{
			Disabler: store.StoreConfig(),
		})
	}

	if ethConfig := mysql.MustNewEthStoreConfigFromViper(); ethConfig.Enabled {
		ctx.ethDB = ethConfig.MustOpenOrCreate(mysql.StoreOption{
			Disabler: store.EthStoreConfig(),
		})
	}

	if redis, ok := redis.MustNewRedisStoreFromViper(store.StoreConfig()); ok {
		ctx.cfxCache = redis
	}

	return ctx
}

func (ctx *storeContext) Close() {
	if ctx.cfxDB != nil {
		ctx.cfxDB.Close()
	}

	if ctx.ethDB != nil {
		ctx.ethDB.Close()
	}

	if ctx.cfxCache != nil {
		ctx.cfxCache.Close()
	}
}

type syncContext struct {
	storeContext
	syncCfx *sdk.Client
	subCfx  *sdk.Client
	syncEth *web3go.Client
}

func mustInitSyncContext(storeCtx storeContext) syncContext {
	sc := syncContext{storeContext: storeCtx}

	if storeCtx.cfxDB != nil || storeCtx.cfxCache != nil {
		sc.syncCfx = rpc.MustNewCfxClientFromViper(rpc.WithClientHookMetrics(true))
		sc.subCfx = rpc.MustNewCfxWsClientFromViper()
	}

	if storeCtx.ethDB != nil {
		sc.syncEth = rpc.MustNewEthClientFromViper(rpc.WithClientHookMetrics(true))
	}

	return sc
}

func (ctx *syncContext) Close() {
	// Usually, storeContext will be defer closed by itself
	// ctx.storeContext.Close()
	if ctx.syncCfx != nil {
		ctx.syncCfx.Close()
	}

	if ctx.subCfx != nil {
		ctx.subCfx.Close()
	}

	// not provided yet!
	// ctx.syncEth.Close()
}
