# Pitfall 3

## Findings Summary

> See the [Per-Finding Index](../../README.md#per-finding-index) for the complete findings table with PoC types, fix status, and similarity values.

| # | Project | PoC Type | Issue Reference |
|---|---------|----------|-----------------|
| 16 | [Sei](https://github.com/sei-protocol/sei-chain) | Mock | [sei-chain#2355](https://github.com/sei-protocol/sei-chain/issues/2355) |
| 17 | [Mitosis](https://github.com/mitosis-org/chain/) | Mock | [chain#87](https://github.com/mitosis-org/chain/issues/87) |
| 18 | [tgrade](https://github.com/confio/tgrade/) | Mock | [abci.go#L36-L50](https://github.com/confio/tgrade/blob/main/x/poe/abci.go#L36-L50) |
| 19 | [Symphony](https://github.com/Orchestra-Labs/symphony) | Mock | [before_send.go#L121-L130](https://github.com/Orchestra-Labs/symphony/blob/main/x/concentrated-liquidity/pool_hooks.go#L121-L130) |
| 20 | [Osmosis](https://github.com/osmosis-labs/osmosis) | Mock | [osmosis#9499](https://github.com/osmosis-labs/osmosis/issues/9499) |
| 21 | [Mantra](https://github.com/MANTRA-Chain/mantrachain) | Mock | [mantrachain#430](https://github.com/MANTRA-Chain/mantrachain/pull/430) |
| 22 | [Neutron](https://github.com/neutron-org/neutron) | Mock | [before_send.go#L165](https://github.com/neutron-org/neutron/blob/daf306ddee9402879acad215dd2e5f2d99f49c8f/x/tokenfactory/keeper/before_send.go#L165) |
| 23 | [Tower](https://github.com/quasar-finance/quasar) | Mock | [before_send.go#L159-L162](https://github.com/quasar-finance/quasar/blob/main/x/tokenfactory/keeper/before_send.go#L159-L162) |
| 24 | [Phoenix](https://github.com/phoenix-directive/core) | Mock | [before_send.go#L138-L146](https://github.com/phoenix-directive/core/blob/release/v2.18/x/tokenfactory/keeper/before_send.go#L138-L146) |
| 25 | [JMES](https://github.com/jmesworld/core) | Mock | [before_send.go#L135-L139](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L135-L139) |

---

<details>
<summary><b>Key Word Match</b></summary>

`/\.WithEventManager/ NOT EventManager().Events() NOT path:handler.go NOT path:context.go NOT path:baseapp.go`

or

`/router\.Handler\(/ NOT EventManager().Events() NOT path:baseapp.go`

or

`.CacheContext() NOT EventManager().Events() NOT path:baseapp.go` and `Cosmos Sdk version > 0.46.3` for chain-registry project

</details>


## PoC

### PoC #19–#25: `WithEventManager` — TokenFactory (Osmosis, Mantra, Neutron, Tower, Phoenix, JMES, Symphony)
x/tokenfactory/keeper/event_in_before_test.go
```go
package keeper_test

import (
	"fmt"
	"os"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/osmosis-labs/osmosis/v30/x/tokenfactory/types"
)

func (s *KeeperTestSuite) TestBeforeSendHookEvent() {
	s.SkipIfWSL()
	for _, tc := range []struct {
		desc     string
		wasmFile string
		sendMsgs []SendMsgTestCase
	}{
		{
			desc:     "74 events when sending 1 of factorydenom, but just 44 events now",
			wasmFile: "./testdata/recall2.wasm",
			sendMsgs: []SendMsgTestCase{
				{
					desc: "74 events when sending 1 of factorydenom, but just 44 events now",
					msg: func(factorydenom string) *banktypes.MsgSend {
						return banktypes.NewMsgSend(
							s.TestAccs[0],
							s.TestAccs[1],
							sdk.NewCoins(sdk.NewInt64Coin(factorydenom, 1)),
						)
					},
					expectPass: true,
				},
			},
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// setup test
			s.SetupTest()
			wasmCode, err := os.ReadFile(tc.wasmFile)
			s.Require().NoError(err, "test1: %v", tc.desc)
			codeID, _, err := s.contractKeeper.Create(s.Ctx, s.TestAccs[0], wasmCode, nil)
			s.Require().NoError(err, "test2: %v", tc.desc)
			cosmwasmAddress, _, err := s.contractKeeper.Instantiate(s.Ctx, codeID, s.TestAccs[0], s.TestAccs[0], []byte("{}"), "", sdk.NewCoins())
			s.Require().NoError(err, "test3: %v", tc.desc)
			res, err := s.msgServer.CreateDenom(s.Ctx, types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
			s.Require().NoError(err, "test: %v", tc.desc)
			denom := res.GetNewTokenDenom()
			_, err = s.msgServer.Mint(s.Ctx, types.NewMsgMint(s.TestAccs[0].String(), sdk.NewInt64Coin(denom, 999999999999999999)))
			s.Require().NoError(err)
			s.FundAcc(sdk.MustAccAddressFromBech32(s.TestAccs[0].String()), sdk.Coins{sdk.NewInt64Coin("foo", 100000000000)})
			_, err = s.msgServer.SetBeforeSendHook(s.Ctx.WithGasMeter(storetypes.NewGasMeter(600000)), types.NewMsgSetBeforeSendHook(s.TestAccs[0].String(), denom, cosmwasmAddress.String()))
			s.Require().NoError(err, "test: %v", tc.desc)
			s.App.BankKeeper.SendCoins(s.Ctx.WithGasMeter(storetypes.NewGasMeter(600000)), s.TestAccs[0], cosmwasmAddress, sdk.NewCoins(sdk.NewInt64Coin(denom, 899999999999999999)))
			denoms, beforeSendHooks := s.App.TokenFactoryKeeper.GetAllBeforeSendHooks(s.Ctx)
			s.Require().Equal(beforeSendHooks, []string{cosmwasmAddress.String()})
			s.Require().Equal(denoms, []string{denom})

			//Halt the chain

			for _, sendTc := range tc.sendMsgs {
				// Create a new context for this send operation to capture events
				sendCtx := s.Ctx.WithGasMeter(storetypes.NewGasMeter(600000))

				// Execute the bank send operation
				s.bankMsgServer.Send(sendCtx, sendTc.msg(denom))

				// Check events in the context
				events := sendCtx.EventManager().Events()
				s.Require().NotEmpty(events, "Expected events to be present in context after send operation")

				// Log all events for debugging
				for i, event := range events {
					fmt.Printf("Event %d: Type=%s, Attributes=%v\n", i, event.Type, event.Attributes)
				}
				s.Require().Equal(len(events), 44, "Expected exactly 44 events") // actual number of events should bee 74
			}
		})
	}
}

```

### PoC #17: `Handler.Router` — Mitosis
chain/x/evmgov/keeper/keeper_test.go
```go
package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/mitosis-org/chain/x/evmgov/keeper"
)

// mockBankMsgServer implements banktypes.MsgServer interface for testing
type mockBankMsgServer struct{}

func (m *mockBankMsgServer) Send(ctx context.Context, msg *banktypes.MsgSend) (*banktypes.MsgSendResponse, error) {
	// Convert context.Context to sdk.Context to access EventManager
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Emit bank transfer event (like real bank module)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			banktypes.EventTypeTransfer,
			sdk.NewAttribute(banktypes.AttributeKeyRecipient, msg.ToAddress),
			sdk.NewAttribute(banktypes.AttributeKeySender, msg.FromAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)
	return &banktypes.MsgSendResponse{}, nil
}

func (m *mockBankMsgServer) MultiSend(ctx context.Context, msg *banktypes.MsgMultiSend) (*banktypes.MsgMultiSendResponse, error) {
	return &banktypes.MsgMultiSendResponse{}, nil
}

func (m *mockBankMsgServer) UpdateParams(ctx context.Context, msg *banktypes.MsgUpdateParams) (*banktypes.MsgUpdateParamsResponse, error) {
	return &banktypes.MsgUpdateParamsResponse{}, nil
}

func (m *mockBankMsgServer) SetSendEnabled(ctx context.Context, msg *banktypes.MsgSetSendEnabled) (*banktypes.MsgSetSendEnabledResponse, error) {
	return &banktypes.MsgSetSendEnabledResponse{}, nil
}

func TestExecuteMessageWithRealHandler(t *testing.T) {
	// Create real message service router
	router := baseapp.NewMsgServiceRouter()

	// Create codec with proper interface registry
	registry := codectypes.NewInterfaceRegistry()
	banktypes.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	// Set up router with interface registry
	router.SetInterfaceRegistry(registry)

	// Register mock bank message service handler
	mockBankServer := &mockBankMsgServer{}
	banktypes.RegisterMsgServer(router, mockBankServer)

	// Create keeper
	k, err := keeper.NewKeeper(cdc, router)
	require.NoError(t, err)

	// Create test context with event manager
	ctx := sdk.Context{}
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// Create a test bank transfer message
	testMsg := &banktypes.MsgSend{
		FromAddress: "cosmos1sender",
		ToAddress:   "cosmos1receiver",
		Amount:      sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1000))),
	}

	// Execute the message
	err = k.ExecuteMessage(ctx, testMsg)

	// Should succeed
	require.NoError(t, err)

	// Check that the bank transfer event was emitted
	events := ctx.EventManager().Events()
	require.Len(t, events, 0, "Expected exactly zero event to be emitted")
}

```

### PoC #16, #18: `CacheContext` — Sei, tgrade


x/epoch/types/hooks_test.go
```golang
package types_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sei-protocol/sei-chain/x/epoch/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

type mockEpochHooksWithEvents struct {
	beforeEpochStartCalled bool
}

func (h *mockEpochHooksWithEvents) AfterEpochEnd(_ sdk.Context, _ types.Epoch) {
	// not implemented for this test
}

func (h *mockEpochHooksWithEvents) BeforeEpochStart(ctx sdk.Context, epoch types.Epoch) {
	h.beforeEpochStartCalled = true

	// Emit an event during BeforeEpochStart
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeNewEpoch,
			sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprint(epoch.CurrentEpoch)),
			sdk.NewAttribute(types.AttributeEpochTime, epoch.CurrentEpochStartTime.String()),
			sdk.NewAttribute(types.AttributeEpochHeight, fmt.Sprint(epoch.CurrentEpochHeight)),
		),
	)
}

func TestBeforeEpochStartEmitsEvents(t *testing.T) {
	hooks := &mockEpochHooksWithEvents{}
	multiHooks := types.MultiEpochHooks{
		hooks,
	}

	db := tmdb.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, tmproto.Header{}, false, nil)

	epoch := types.Epoch{
		CurrentEpoch:          1,
		CurrentEpochHeight:    100,
		CurrentEpochStartTime: ctx.BlockTime(),
	}

	// Check initial event count
	initialEventCount := len(ctx.EventManager().Events())

	// Call BeforeEpochStart which should emit an event
	multiHooks.BeforeEpochStart(ctx, epoch)

	// Verify the hook was called
	require.True(t, hooks.beforeEpochStartCalled)

	// Check that exactly one event was emitted
	finalEventCount := len(ctx.EventManager().Events())
	require.Equal(t, initialEventCount, finalEventCount, "Expected exactly one event to be emitted")
}
```