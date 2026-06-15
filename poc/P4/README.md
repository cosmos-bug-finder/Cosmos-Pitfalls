# Pitfall 4

## Findings Summary

> See the [Per-Finding Index](../../README.md#per-finding-index) for the complete findings table with PoC types, fix status, and similarity values.

| # | Project | PoC Type | Issue Reference |
|---|---------|----------|-----------------|
| 26 | [Stride](https://github.com/Stride-Labs/stride) | Mock | [stride#1402](https://github.com/Stride-Labs/stride/issues/1402) |
| 27 | [SSC](https://github.com/sagaxyz/ssc) | Mock | [ssc#38](https://github.com/sagaxyz/ssc/issues/38) |
| 28 | [UptickNetwork](https://github.com/UptickNetwork/uptick) | Mock | [uptick#27](https://github.com/UptickNetwork/uptick/issues/27) |
| 29 | [K***](https://github.com/K***-Labs/K***) | Mock | [K***#2077](https://github.com/K***-Labs/K***/issues/2077) |
| 30 | [Irisnet](https://github.com/irisnet/irishub) | Mock | [irishub#2997](https://github.com/irisnet/irishub/issues/2997) |
| 31 | [XteLabs](https://github.com/xtelabs/xtechain) | Mock | [xtechain#261](https://github.com/xtelabs/xtechain/issues/261) |
| 32 | [Planq](https://github.com/planq-network/planq) | Mock | [planq#282](https://github.com/planq-network/planq/issues/282) |
| 33 | [Mezod](https://github.com/mezo-org/mezod) | Mock | [mezod#536](https://github.com/mezo-org/mezod/issues/536) |
| 34 | [Tenet-Evmos](https://github.com/tenet-org/tenet-evmos) | Mock | [tenet-evmos#323](https://github.com/tenet-org/tenet-evmos/issues/323) |
| 35 | [MTT Chain](https://github.com/mtt-labs/mtt-chain) | Mock | [mtt-chain#6](https://github.com/mtt-labs/mtt-chain/issues/6) |
| 36 | [Hetu Chain](https://github.com/hetu-project/hetu-chain) | Mock | [hetu-chain#2](https://github.com/hetu-project/hetu-chain/issues/2) |
| 37 | [H*** Core](https://github.com/H***-network/H***-core) | Mock | [H***-core#36](https://github.com/H***-network/H***-core/issues/36) |
| 38 | [LambdaVM](https://github.com/LambdaIM/lambdavm) | Mock | [lambdavm#255](https://github.com/LambdaIM/lambdavm/issues/255) |
| 39 | [G***](https://github.com/GPTx-global/G***) | Mock | [G***-v1#32](https://github.com/G***finglobal/G***-v1/issues/32) |
| 40 | [Bridgeless Core](https://github.com/hyle-team/bridgeless-core) | Mock | [bridgeless-core#99](https://github.com/Bridgeless-Project/bridgeless-core/issues/99) |
| 41 | [Egochain](https://github.com/EgorasMarket/Egochain-Blockchain) | Mock | [Egochain#223](https://github.com/EgorasMarket/Egochain-Blockchain/issues/223) |
| 42 | [CVN](https://github.com/cvn-network/cvn) | Mock | [cvn#3](https://github.com/cvn-network/cvn/issues/3) |
| 43 | [DE-EVM](https://github.com/depaasecology/de-evm) | Mock | [de-evm#1](https://github.com/depaasecology/de-evm/issues/1) |
| 44 | [Rarimo Core](https://github.com/rarimo/rarimo-core) | Mock | [rarimo-core#44](https://github.com/rarimo/rarimo-core/issues/44) |
| 45 | [Cosmos EVM](https://github.com/cosmos/evm) | Mock | [evm#628](https://github.com/cosmos/evm/issues/628) |
| 46 | [BRC Chain](https://github.com/brcchain/brcchain) | Mock | [brcchain#1](https://github.com/brcchain/brcchain/issues/1) |
| 47 | [Rollup EVM](https://github.com/airchains-network/rollup-evm) | Mock | [rollup-evm#3](https://github.com/airchains-network/rollup-evm/issues/3) |
| 48 | [Artela](https://github.com/artela-network/artela) | Mock | [artela#191](https://github.com/artela-network/artela/issues/191) |
| 49 | [OLLO](https://github.com/OllO-Station/ollo) | Mock | [ollo#75](https://github.com/OllO-Station/ollo/issues/75) |

---

<details>
<summary><b>Key Word Match</b></summary>

`.CacheContext() AND EventManager().Events() NOT path:baseapp.go`

and `Cosmos Sdk version > 0.46.3`

</details>

## PoC
### PoC #27: Saga (SSC)
ssc/x/epochs/types/hooks_test.go
```go
package types

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyFuncIfNoError(t *testing.T) {
	storeKey := storetypes.NewKVStoreKey("test")
	memStoreKey := storetypes.NewMemoryStoreKey("test_mem")

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	t.Run("number text for event", func(t *testing.T) {
		f := func(ctx sdk.Context) error {
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				"test_event",
				sdk.NewAttribute("key", "value"),
			))
			return nil
		}

		err := applyFuncIfNoError(ctx, f)
		assert.NoError(t, err)
		finalEventCount := len(ctx.EventManager().Events())
		assert.Equal(t, 2, finalEventCount)
	})
}


```


### PoC #29–#49: EVM (ethermint fork cluster)

Here, for evmos and uptick we provide the use of cachecontext to validate that the event emit has been emitted internally.  We simplify the validation to make the problem clearer.

Meanwhile, we can clearly check the use case of P4 in the following problem code.

https://github.com/K***-Labs/ethermint/blob/12fdf9416415a8ec0a96df5a42e468b0e072f5cc/x/evm/keeper/state_transition.go#L232-L235


/x/evm/keeper/cache_applyposttx_test.go
```go

package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createTestContext(t *testing.T) sdk.Context {
	t.Helper()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	ctx := sdk.NewContext(ms, tmproto.Header{
		Height: 1,
		Time:   time.Now(),
	}, false, nil)

	return ctx
}

func PostTxProcessing(ctx sdk.Context, fromAddr, toAddr common.Address, value *big.Int, txHash common.Hash) error {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"transaction_completed",
			sdk.NewAttribute("tx_hash", txHash.Hex()),
			sdk.NewAttribute("status", "success"),
		),
	)

	return nil
}

func mockApplyTransaction(ctx sdk.Context, tx *ethtypes.Transaction) error {
	tmpCtx, commit := ctx.CacheContext()

	fromAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	toAddr := common.HexToAddress("0x0987654321098765432109876543210987654321")

	if err := PostTxProcessing(tmpCtx, fromAddr, toAddr, tx.Value(), tx.Hash()); err != nil {
		return err
	}

	if commit != nil {
		commit()
		ctx.EventManager().EmitEvents(tmpCtx.EventManager().Events())
	}

	return nil
}

func TestMockApplyTransaction(t *testing.T) {
	ctx := createTestContext(t)
	tx := ethtypes.NewTransaction(
		0,
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		big.NewInt(1000000000000000000),
		21000,
		big.NewInt(20000000000),
		[]byte{},
	)
	err := mockApplyTransaction(ctx, tx)
	require.NoError(t, err)
	eventCount := len(ctx.EventManager().Events())
	require.Equal(t, 2, eventCount)
}


```

### PoC #26: Stride
stride/stride/utils/cache_ctx_test.go
```go
package utils_test

import (
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Stride-Labs/stride/v24/utils"
)

func (s *UtilsTestSuite) TestCacheCtxEmitEvent() {
	ctx := s.Ctx.WithGasMeter(types.NewGasMeter(10_000))
	utils.ApplyFuncIfNoError(ctx, func(c sdk.Context) error {
		c.EventManager().EmitEvent(sdk.NewEvent("test", sdk.NewAttribute("test", "test")))
		return nil
	})
	s.Equal(2, len(ctx.EventManager().Events()))
}

```

### PoC #28: UptickNetwork
Here, for evmos and uptick we provide the use of cachecontext to validate that the event emit has been emitted internally.  We simplify the validation to make the problem clearer.

Meanwhile, we can clearly check the use case of P4 in the following problem code.

https://github.com/UptickNetwork/uptick/blob/a062c96c63bf6b3f9d5db6251b11ad78f69fe978/x/evmIBC/keeper/ibc_hook.go#L71-L72

/x/evmIBC/keeper/cache_test.go
```go
package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMockRecvPacket2(t *testing.T) {
	ctx := sdk.Context{}
	cctx, write := ctx.CacheContext()
	event := sdk.NewEvent(
		"convert_coin",
		sdk.NewAttribute("test", "value"),
	)
	cctx.EventManager().EmitEvent(event)
	
	//test for dup emit process.
	write()
	ctx.EventManager().EmitEvents(cctx.EventManager().Events())
	events := ctx.EventManager().Events()
	require.Len(t, events, 2, "Context should have exactly 2 events")
}
```

### PoC #49: OllO
x/exchange/abci_poc_test.go
```golang
package exchange_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const poc_evtType = "exchange_batch_event"

// mockMidBlocker mirrors x/exchange/abci.go MidBlocker (lines 32-48): same market
// loop + cacheCtx/writeCache + manual EmitEvents structure, so the double-emit is
// reproduced exactly as it happens on-chain.
func mockMidBlocker(ctx sdk.Context) {
	markets := []string{"market-1"} // stand-in for k.IterateAllMarkets(ctx, ...)
	for range markets {
		cacheCtx, writeCache := ctx.CacheContext()
		func() {
			defer func() {
				if r := recover(); r != nil {
					_ = r // mirrors k.Logger(ctx).Error("panic in batch matching", ...)
				}
			}()
			if err := mockRunBatchMatching(cacheCtx); err != nil {
				_ = err // mirrors k.Logger(ctx).Error("failed to run batch matching", ...)
			} else {
				writeCache()
				// Pass events emitted inside cached context to the original context.
				ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
			}
		}()
	}
}

// mockRunBatchMatching stands in for k.RunBatchMatching: it emits an event into
// the cached context, just like real batch matching does.
func mockRunBatchMatching(ctx sdk.Context) error {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(poc_evtType, sdk.NewAttribute("orderer", "alice")),
	)
	return nil
}

func TestExchangeMidBlockerDoubleEmit_PoC(t *testing.T) {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(storetypes.NewKVStoreKey("poc"), storetypes.StoreTypeIAVL, db)
	require.NoError(t, cms.LoadLatestVersion())
	ctx := sdk.NewContext(cms, tmproto.Header{Height: 1}, false, log.NewNopLogger())

	mockMidBlocker(ctx)

	var count int
	for _, ev := range ctx.EventManager().Events() {
		if ev.Type == poc_evtType {
			count++
		}
	}
	require.Equal(t, 2, count)
}
```
