# Pitfall 5

## Findings Summary

> See the [Per-Finding Index](../../README.md#per-finding-index) for the complete findings table with PoC types, fix status, and similarity values.

| # | Project | PoC Type | Issue Reference |
|---|---------|----------|-----------------|
| 50 | [Evmos](https://github.com/evmos/evmos) | Full Chain | CVE-2021-25837 |

---

<details>
<summary><b>Key Word Match</b></summary>


`/type Keeper struct \{[\s\S]*Cache[\s\S]*?\}[\s\S]*NewKeeper/ AND cosmos AND path:keeper.go`

We collected all the [keepers](../../record/record/keeper/result/keeper_structs_filtered.json) of the projects that were registered in the chain-registry by weeding out the common fields. 

We manually checked all the uncommon fields, the only discovery was `filecache` under [`Band`](https://github.com/bandprotocol/chain/blob/master/pkg/filecache/filecache.go) project, we collected more projects hosted in github but not appearing in chain-registry using cache as a keyword, and didn't get a valid discovery. It is worth noting that although P5 exists in the band project, we found that the issue cannot be exploited because accessing the filecache stored in the keeper requires an index stored in the Context, which is normally rolled back when a transaction fails resulting in no access to the dirty-written filecache.

</details>

## PoC
### PoC #50: Evmos

see [evmos-case](../../high-light-findings/evm-ethermint)

### Band (excluded — not exploitable)

/chain/x/oracle/keeper/msg_server_test.go
```golang
package keeper_test

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/chain/v3/x/oracle/keeper"
	"github.com/bandprotocol/chain/v3/x/oracle/types"
)

// mockRunTx simulates transaction execution with cache context for multiple messages
func (s *KeeperTestSuite) mockRunTx(msgs ...sdk.Msg) error {
	cacheCtx, commit := s.ctx.CacheContext()
	msgServer := keeper.NewMsgServerImpl(s.oracleKeeper)

	// Execute all messages in cache context
	for _, msg := range msgs {
		switch m := msg.(type) {
		case *types.MsgCreateDataSource:
			_, err := msgServer.CreateDataSource(cacheCtx, m)
			if err != nil {
				return err
			}
		default:
			return types.ErrBadDrbgInitialization.Wrap("unsupported message type")
		}
	}

	// If all messages succeed, commit the changes
	commit()
	return nil
}

func (s *KeeperTestSuite) TestCreateDataSource() {
	initialDataSourceCount := s.oracleKeeper.GetDataSourceCount(s.ctx)
	executable := []byte("test executable content")
	fee := sdk.NewCoins(sdk.NewCoin("uband", math.NewInt(1000)))
	msg1 := &types.MsgCreateDataSource{
		Name:        "Test Data Source 1",
		Description: "First test data source",
		Executable:  executable,
		Fee:         fee,
		Treasury:    treasury.String(),
		Owner:       owner.String(),
		Sender:      owner.String(),
	}

	msg2 := &types.MsgEditDataSource{
		Name:        "Test Data Source 2",
		Description: "Second test data source",
		Executable:  executable,
		Fee:         fee,
		Treasury:    treasury.String(),
		Owner:       alice.String(),
		Sender:      alice.String(),
	}

	// 1.Execute both messages using mock transaction, which must fail.
	s.mockRunTx(msg1, msg2)
	// 2.KvStorage is not been modified.
	// Verify final state after transaction
	finalDataSourceCount := s.oracleKeeper.GetDataSourceCount(s.ctx)
	// KvStorage is not been modified.
	s.Require().Equal(initialDataSourceCount, finalDataSourceCount, "Two data sources should be added")
	// 3.Verify file cache contains the executable files by their SHA256-based filenames : 338bee2f3f25945c43183f82de271d2b7d44c50e5f63d47a63892a2e353e2cc3
	a := s.oracleKeeper.GetFile("338bee2f3f25945c43183f82de271d2b7d44c50e5f63d47a63892a2e353e2cc3")
	s.Require().NotNil(a)
}


```
