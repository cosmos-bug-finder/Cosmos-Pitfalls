# Cosmos Application Development Pitfall

A comprehensive security research study identifying common vulnerabilities and implementation issues across Cosmos ecosystem blockchain projects. This research identifies and documents security pitfalls and vulnerabilities found across multiple Cosmos-based blockchain projects. The findings span 8 different categories (P1-P8). This repo include the PoC and key info about our findings.

## A. Highlighted Findings

The following shows what we consider to be high-impact findings in this study.

### Asset Stolen

- [Evmos/ethermint](./high-light-findings/evm-ethermint) - `Evmos/ethermint` is a widely used evm module in Cosmos ecosystem, including Cronos($3.13B), Injective($310M), Mantra($71.84M), Kava($61M), Zeta($71M) and so on.([CVE](https://www.cve.org/CVERecord?id=CVE-2021-25837))

- [Secret](./high-light-findings/secret) - Prevented the potential $50 million value assets theft ([Acknowledgment](./record/acknowledgments/secret.md))

### Chain Halt

- [Mantra](./high-light-findings/mantra-tokenfactory) - Market capitalization of $50 million([CVE](https://github.com/MANTRA-Chain/mantrachain/security/advisories/GHSA-qwvm-wqq8-8j69))
- [Injective](./high-light-findings/Injective) - Market capitalization of $300 million ([Acknowledgment](./record/acknowledgments/injective.md))

## B. Findings Summary

- ✅ **Fixed**: A fix has been deployed in a release or merged commit.
- **Acknowledged**: Project team acknowledged and committed to fixing, but no public fix was available at the time of writing.
- **-**: No response received from the project team.

> **Note on "Privileged Account"**: Entries marked with "Privileged Account" in the Deploy Permission column indicate that the target chain requires a governance proposal or whitelisted account to deploy CosmWasm contracts. This restricts the attack surface to authorized deployers only.

### Per-Finding Index

| # | Pitfall | Project | PoC Type | Issue Reference | Deploy Permission | Inherited Module | Fix | Credit |
|---|---------|---------|----------|-----------------|--------------------|-----------|-----|--------|
| 1 | **P1** | [Mantra](https://github.com/MANTRA-Chain/mantrachain) | UnitTest | [cosmos-sdk#25303](https://github.com/cosmos/cosmos-sdk/pull/25303) | - | tokenfactory | ✅ [Fixed](https://github.com/MANTRA-Chain/mantrachain/pull/438/changes) | [Issue Credit](https://github.com/cosmos/cosmos-sdk/pull/25303) |
| 2 |  | [Osmosis](https://github.com/osmosis-labs/osmosis) | UnitTest | [osmosis#9511](https://github.com/osmosis-labs/osmosis/pull/9511) | [Privileged Account](https://docs.osmosis.zone/cosmwasm/local/submit-wasm-proposal) | tokenfactory  | ✅ [Fixed](https://github.com/osmosis-labs/osmosis/pull/9511) | [Issue Credit](https://github.com/osmosis-labs/osmosis/pull/9511) |
| 3 |  | [Neutron](https://github.com/neutron-org/neutron) | UnitTest | [neutron#978](https://github.com/neutron-org/neutron/pull/978) | [Privileged Account](https://github.com/Orchestra-Labs/symphony/blob/main/x/tokenfactory/keeper/msg_server.go#L210-L217) | tokenfactory | [Acknowledged](https://github.com/neutron-org/neutron/pull/978) | - |
| 4 |  | [Tower](https://github.com/quasar-finance/quasar) | UnitTest | [before_send.go#L161](https://github.com/quasar-finance/quasar/blob/main/x/tokenfactory/keeper/before_send.go#L161) | [Privileged Account](https://github.com/quasar-finance/quasar/blob/main/demos/upgrade-handler/v3.0.0/post_upgrade.sh) | tokenfactory | - | - |
| 5 |  | [Phoenix](https://github.com/phoenix-directive/core) | UnitTest | [before_send.go#L144](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L144) | - | tokenfactory | - | - |
| 6 |  | [JMES](https://github.com/jmesworld/core) | UnitTest | [before_send.go#L144](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L144) | - | tokenfactory | - | - |
| 7 |  | [Symphony](https://github.com/Orchestra-Labs/symphony) | UnitTest | [before_send.go#L165](https://github.com/Orchestra-Labs/symphony/blob/main/x/tokenfactory/keeper/before_send.go#L165) | [Privileged Account](https://github.com/Orchestra-Labs/symphony/blob/main/x/cosmwasmpool/README.md#governance-and-code-id-management) | tokenfactory | - | - |
| 8 |  | [Juno](https://github.com/CosmosContracts/juno) | UnitTest | [contracts.go#L70-L87](https://github.com/CosmosContracts/juno/blob/main/x/cw-hooks/keeper/contracts.go#L70-L87) | - |  | [Acknowledged](https://github.com/CosmosContracts/juno/commit/7dbbf115e2c0c33b4117fb26c1e11e6e1a324b67) | - |
| 9 | **P2** | [Mantra](https://github.com/MANTRA-Chain/mantrachain) | UnitTest | [cosmos-sdk#25303](https://github.com/cosmos/cosmos-sdk/pull/25303) | - | tokenfactory | ✅ [Fixed](https://github.com/MANTRA-Chain/mantrachain/pull/438/changes) | [CVE-2025-61595](https://github.com/MANTRA-Chain/mantrachain/security/advisories/GHSA-qwvm-wqq8-8j69) |
| 10 |  | [Neutron](https://github.com/neutron-org/neutron) | UnitTest | [neutron#978](https://github.com/neutron-org/neutron/pull/978) | [Privileged Account](https://github.com/Orchestra-Labs/symphony/blob/main/x/tokenfactory/keeper/msg_server.go#L210-L217) | tokenfactory | [Acknowledged](https://github.com/neutron-org/neutron/pull/978) | [Private Credit](./record/acknowledgments/neutron.md) |
| 11 |  | [Tower](https://github.com/quasar-finance/quasar) | UnitTest | [before_send.go#L161](https://github.com/quasar-finance/quasar/blob/main/x/tokenfactory/keeper/before_send.go#L161) | [Privileged Account](https://github.com/quasar-finance/quasar/blob/main/demos/upgrade-handler/v3.0.0/post_upgrade.sh) | tokenfactory | - | - |
| 12 |  | [Phoenix](https://github.com/phoenix-directive/core) | UnitTest | [before_send.go#L144](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L144) | - | tokenfactory | - | - |
| 13 |  | [JMES](https://github.com/jmesworld/core) | UnitTest | [before_send.go#L144](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L144) | - | tokenfactory | - | - |
| 14 |  | [Symphony](https://github.com/Orchestra-Labs/symphony) | UnitTest | [before_send.go#L165](https://github.com/Orchestra-Labs/symphony/blob/main/x/tokenfactory/keeper/before_send.go#L165) | [Privileged Account](https://github.com/Orchestra-Labs/symphony/blob/main/x/cosmwasmpool/README.md#governance-and-code-id-management) | tokenfactory | - | - |
| 15 |  | [Juno](https://github.com/CosmosContracts/juno) | Mock | [contracts.go#L70-L87](https://github.com/CosmosContracts/juno/blob/main/x/cw-hooks/keeper/contracts.go#L70-L87) | - |  | [Acknowledged](https://github.com/CosmosContracts/juno/commit/7dbbf115e2c0c33b4117fb26c1e11e6e1a324b67) | [Private Credit](./record/acknowledgments/juno.md) |
| 16 | **P3** | [Sei](https://github.com/sei-protocol/sei-chain) | Mock | [sei-chain#2355](https://github.com/sei-protocol/sei-chain/issues/2355) | - |  | - | - |
| 17 |  | [Mitosis](https://github.com/mitosis-org/chain/) | Mock | [chain#87](https://github.com/mitosis-org/chain/issues/87) | - |  | - | - |
| 18 |  | [tgrade](https://github.com/confio/tgrade/) | Mock | [abci.go#L36-L50](https://github.com/confio/tgrade/blob/main/x/poe/abci.go#L36-L50) | - |  | - | - |
| 19 |  | [Symphony](https://github.com/Orchestra-Labs/symphony) | Mock | [pool_hooks.go#L121-L130](https://github.com/Orchestra-Labs/symphony/blob/main/x/concentrated-liquidity/pool_hooks.go#L121-L130) | - |  | - | - |
| 20 |  | [Osmosis](https://github.com/osmosis-labs/osmosis) | Mock | [osmosis#9499](https://github.com/osmosis-labs/osmosis/issues/9499) | - | tokenfactory  | - | - |
| 21 |  | [Mantra](https://github.com/MANTRA-Chain/mantrachain) | Mock | [mantrachain#430](https://github.com/MANTRA-Chain/mantrachain/pull/430) | - | tokenfactory | ✅ [Fixed](https://github.com/MANTRA-Chain/mantrachain/pull/430) | [Issue Credit](https://github.com/MANTRA-Chain/mantrachain/issues/427) |
| 22 |  | [Neutron](https://github.com/neutron-org/neutron) | Mock | [before_send.go#L165](https://github.com/neutron-org/neutron/blob/daf306ddee9402879acad215dd2e5f2d99f49c8f/x/tokenfactory/keeper/before_send.go#L165) | - | tokenfactory | [Acknowledged](https://github.com/neutron-org/neutron/commit/3e26ebbb8f3aed251af2828aca9d29666feb8622) | - |
| 23 |  | [Tower](https://github.com/quasar-finance/quasar) | Mock | [before_send.go#L159-L162](https://github.com/quasar-finance/quasar/blob/main/x/tokenfactory/keeper/before_send.go#L159-L162) | - | tokenfactory | - | - |
| 24 |  | [Phoenix](https://github.com/phoenix-directive/core) | Mock | [before_send.go#L138-L146](https://github.com/phoenix-directive/core/blob/release/v2.18/x/tokenfactory/keeper/before_send.go#L138-L146) | - | tokenfactory | - | - |
| 25 |  | [JMES](https://github.com/jmesworld/core) | Mock | [before_send.go#L135-L139](https://github.com/jmesworld/core/blob/main/x/tokenfactory/keeper/before_send.go#L135-L139) | - | tokenfactory | - | - |
| 26 | **P4** | [Stride](https://github.com/Stride-Labs/stride) | Mock | [stride#1402](https://github.com/Stride-Labs/stride/issues/1402) | - |  | ✅ [Fixed](https://github.com/Stride-Labs/stride/pull/1415) | [Issue Credit](https://github.com/Stride-Labs/stride/pull/1415) |
| 27 |  | [SSC](https://github.com/sagaxyz/ssc) | Mock | [ssc#38](https://github.com/sagaxyz/ssc/issues/38) | - |  | - | - |
| 28 |  | [UptickNetwork](https://github.com/UptickNetwork/uptick) | Mock | [uptick#27](https://github.com/UptickNetwork/uptick/issues/27) | - |  | ✅ [Fixed](https://github.com/UptickNetwork/uptick/issues/27) | [Release Credit](https://github.com/UptickNetwork/uptick/releases/tag/v0.3.1) |
| 29 |  | [Kava](https://github.com/Kava-Labs/kava) | Mock | [kava#2077](https://github.com/Kava-Labs/kava/issues/2077) | - | EVM | [Acknowledged](https://github.com/Kava-Labs/kava/issues/2077) | [Private Credit](./record/acknowledgments/kava.md) |
| 30 |  | [Irisnet](https://github.com/irisnet/irishub) | Mock | [irishub#2997](https://github.com/irisnet/irishub/issues/2997) | - | EVM | ✅ [Fixed](https://github.com/irisnet/irishub/pull/2998) | [Issue Credit](https://github.com/irisnet/irishub/issues/2997) |
| 31 |  | [XteLabs](https://github.com/xtelabs/xtechain) | Mock | [xtechain#261](https://github.com/xtelabs/xtechain/issues/261) | - | EVM | - | - |
| 32 |  | [Planq](https://github.com/planq-network/planq) | Mock | [planq#282](https://github.com/planq-network/planq/issues/282) | - | EVM | - | - |
| 33 |  | [Mezod](https://github.com/mezo-org/mezod) | Mock | [mezod#536](https://github.com/mezo-org/mezod/issues/536) | - | EVM | - | - |
| 34 |  | [Tenet-Evmos](https://github.com/tenet-org/tenet-evmos) | Mock | [tenet-evmos#323](https://github.com/tenet-org/tenet-evmos/issues/323) | - | EVM | - | - |
| 35 |  | [MTT Chain](https://github.com/mtt-labs/mtt-chain) | Mock | [mtt-chain#6](https://github.com/mtt-labs/mtt-chain/issues/6) | - | EVM | - | - |
| 36 |  | [Hetu Chain](https://github.com/hetu-project/hetu-chain) | Mock | [hetu-chain#2](https://github.com/hetu-project/hetu-chain/issues/2) | - | EVM | ✅ [Fixed](https://github.com/hetu-project/hetu-chain/pull/3) | [Issue Credit](https://github.com/hetu-project/hetu-chain/issues/2) |
| 37 |  | [Helios Core](https://github.com/helios-network/helios-core) | Mock | [helios-core#36](https://github.com/helios-network/helios-core/issues/36) | - | EVM | [Acknowledged](https://github.com/helios-network/helios-core/issues/36) | [Issue Credit](https://github.com/helios-network/helios-core/issues/36) |
| 38 |  | [LambdaVM](https://github.com/LambdaIM/lambdavm) | Mock | [lambdavm#255](https://github.com/LambdaIM/lambdavm/issues/255) | - | EVM | - | - |
| 39 |  | [Guru](https://github.com/GPTx-global/guru) | Mock | [guru-v1#32](https://github.com/gurufinglobal/guru-v1/issues/32) | - | EVM | [Acknowledged](https://github.com/gurufinglobal/guru-v1/issues/32) | [Issue Credit](https://github.com/gurufinglobal/guru-v1/issues/32) |
| 40 |  | [Bridgeless Core](https://github.com/hyle-team/bridgeless-core) | Mock | [bridgeless-core#99](https://github.com/Bridgeless-Project/bridgeless-core/issues/99) | - | EVM | ✅ [Fixed](https://github.com/Bridgeless-Project/bridgeless-core/pull/101) | [Issue Credit](https://github.com/Bridgeless-Project/bridgeless-core/issues/99) |
| 41 |  | [Egochain](https://github.com/EgorasMarket/Egochain-Blockchain) | Mock | [Egochain#223](https://github.com/EgorasMarket/Egochain-Blockchain/issues/223) | - | EVM | - | - |
| 42 |  | [CVN](https://github.com/cvn-network/cvn) | Mock | [cvn#3](https://github.com/cvn-network/cvn/issues/3) | - | EVM | - | - |
| 43 |  | [DE-EVM](https://github.com/depaasecology/de-evm) | Mock | [de-evm#1](https://github.com/depaasecology/de-evm/issues/1) | - | EVM | - | - |
| 44 |  | [Rarimo Core](https://github.com/rarimo/rarimo-core) | Mock | [rarimo-core#44](https://github.com/rarimo/rarimo-core/issues/44) | - | EVM | - | - |
| 45 |  | [Cosmos EVM](https://github.com/cosmos/evm) | Mock | [evm#628](https://github.com/cosmos/evm/issues/628) | - | EVM | ✅ [Fixed](https://github.com/cosmos/evm/commit/7a383702e62b47f1bab82f27d54f6b2f832e41ad) | [Issue Credit](https://github.com/cosmos/evm/issues/628) |
| 46 |  | [BRC Chain](https://github.com/brcchain/brcchain) | Mock | [brcchain#1](https://github.com/brcchain/brcchain/issues/1) | - | EVM | - | - |
| 47 |  | [Rollup EVM](https://github.com/airchains-network/rollup-evm) | Mock | [rollup-evm#3](https://github.com/airchains-network/rollup-evm/issues/3) | - | EVM | - | - |
| 48 |  | [Artela](https://github.com/artela-network/artela) | Mock | [artela#191](https://github.com/artela-network/artela/issues/191) | - | EVM | - | - |
| 49 |  | [OLLO](https://github.com/OllO-Station/ollo) | Mock | [ollo#75](https://github.com/OllO-Station/ollo/issues/75) | - | EVM | - | - |
| 50 | **P5** | [Evmos](https://github.com/evmos/evmos) | Full Chain | [CVE-2021-25837](https://www.cve.org/CVERecord?id=CVE-2021-25837) | - |  | ✅ Fixed | [**CVE-2021-25837**](https://www.cve.org/CVERecord?id=CVE-2021-25837) |
| 51 | **P6** | [Jackal](https://github.com/JackalLabs/canine-chain) | UnitTest | [canine-chain#8](https://github.com/JackalLabs/canine-chain/issues/8) | - | RNS | ✅ [Fixed](https://github.com/JackalLabs/canine-chain/pull/9) | [Issue Credit](https://github.com/JackalLabs/canine-chain/issues/8) |
| 52 |  | [Ignite](https://github.com/ignite/cli/blob/main/docs/versioned_docs/version-v0.26/02-guide) | UnitTest | [cli#2828](https://github.com/ignite/cli/issues/2828) | - | RNS & Loan | ✅ [Fixed](https://github.com/ignite/cli/issues/2828) | [Issue Credit](https://github.com/ignite/cli/issues/2828) |
| 53 |  | [Side Protocol](https://github.com/sideprotocol/ibcswap) | UnitTest | [ibcswap#8](https://github.com/sideprotocol/ibcswap/issues/8) | - |  | ✅ [Fixed](https://github.com/sideprotocol/ibcswap/commit/85f9d32099ebd6558486610bdb9f7918282a590f) | [Issue Credit](https://github.com/sideprotocol/ibcswap/issues/8) |
| 54 |  | [OLLO](https://github.com/OllO-Station/ollo) | UnitTest | [ollo#20](https://github.com/OllO-Station/ollo/issues/20) | - | Loan | ✅ [Fixed](https://github.com/OllO-Station/ollo/commit/9ed0245b09218858a3806064af683cb18e1e1a39) | [Issue Credit](https://github.com/OllO-Station/ollo/issues/20) |
| 55 |  | [Ununifi](https://github.com/exoralayer/chain) | UnitTest | [chain#612](https://github.com/exoralayer/chain/issues/612) | - |  | ✅ [Fixed](https://github.com/exoralayer/chain/commit/eac956e480e264341c13530094c153b11e9fc3a7) | [Issue Credit](https://github.com/exoralayer/chain/issues/612) |
| 56 |  | [Stateset](https://github.com/stateset/core) | UnitTest | [core#1](https://github.com/stateset/core/issues/1) | - |  | - | - |
| 57 |  | [Haven](https://github.com/onomyprotocol/haven) | UnitTest | [haven#9](https://github.com/onomyprotocol/haven/issues/9) | - |  | - | - |
| 58 |  | [ICPlaza](https://github.com/ICPLAZA-org/icplaza) | UnitTest | [icplaza#1](https://github.com/ICPLAZA-org/icplaza/issues/1) | - |  | - | - |
| 59 |  | [hero](https://github.com/strangelove-ventures/hero) | UnitTest | [msg_server_burn.go#L41](https://github.com/strangelove-ventures/hero/blob/d6737b82c5ff4976f9a9a57d75a73d61f84bf395/x/tokenfactory/keeper/msg_server_burn.go#L41) | - |  | - | - |
| 60 |  | [decoGit](https://github.com/ritajeong/decoGit) | UnitTest | [msg_server_buy_sticker.go#L39](https://github.com/ritajeong/decoGit/blob/9cc1398419dfd0dc26577a7709bcbbf919c1d0d2/server/chain/x/decogit/keeper/msg_server_buy_sticker.go#L39) | - |  | - | - |
| 61 |  | [HalbornSecurity](https://github.com/HalbornSecurity/CTFs) | UnitTest | [msg_server_mint_hal.go#L55](https://github.com/HalbornSecurity/CTFs/blob/684f1af02132d4cfa4c6ac04924d8c8391a6e9cf/HalbornCTF_Golang_Cosmos/x/hal/keeper/msg_server_mint_hal.go#L55) | - |  | - | - |
| 62 | **P7** | [Secret Network](https://github.com/scrtlabs/SecretNetwork) | Full Chain | [Ack](./record/acknowledgments/secret.md) | - |  | ✅ [Fixed](./record/acknowledgments/secret.md) | [Bounty Evident](https://www.mintscan.io/secret/address/secret10e2s2ygx7zmdcrnmwphuljl8sg0ytl6aam545t) |
| 63 | **P8** | [Injective](https://github.com/InjectiveFoundation/injective-core) | Full Chain | [app.go#L349-L353](https://github.com/InjectiveFoundation/injective-core/blob/release/v1.18.x/injective-chain/app/app.go#L349-L353) | [Privileged Account](https://docs.injective.network/developers-cosmwasm/smart-contracts/your-first-smart-contract#upload-the-wasm-contract) | Evm & Wasm  | ✅ Fixed | [Private Credit](./record/acknowledgments/injective.md) |
| 64 |  | [PellNetwork](https://github.com/0xPellNetwork/aegis) | Full Chain | [app.go#L331-L341](https://github.com/0xPellNetwork/aegis/blob/main/app/app.go#L331-L341) | - | Evm & Wasm | - | - |
| 65 |  | [Uptick](https://github.com/UptickNetwork/uptick) | Full Chain | [app.go#L89-L96](https://github.com/UptickNetwork/uptick/blob/main/app/app.go#L89-L96) | - | Evm & Wasm | ✅ [Fixed](https://github.com/UptickNetwork/uptick/releases/tag/v0.3.1) | [Release Credit](https://github.com/UptickNetwork/uptick/releases/tag/v0.3.1) |

The vast majority of the projects listed above remain available when we identified the issue (verified by the accessibility of their official websites at the time the issues were identified). In addition, other items have been retained in the list due to their unique significance. For instance, `HalbornSecurity` was discovered as a CTF challenge created by an external party based on one of our previous findings. While `Cosmos/Evm` and `Ignite` are not projects, they are widely utilized as dependency libraries or tutorial code.

### Originality Statement

Regarding prior security work on `P1`–`P8`, to the best of our knowledge, only `P1` and `P8` had been publicly documented before. P1 was previously described in the CosmWasm advisory [CWA-2024-008](https://github.com/CosmWasm/advisories/blob/c3d9c3021ca262b91b1f7f326e2b4cd2d3e5835c/CWAs/CWA-2024-008.md), which we became aware of only after independently identifying the same error pattern during our later-stage collection process. P8 was inspired by a prior blog from [Jump Crypto](https://jumpcrypto.com/resources/stealing-gas-bypassing-ethermint-ante-handlers), which discussed conflicts between the `group/evmos` and `wasmd/evmos` modules, but did not systematically study similar issues across the ecosystem.

To the best of our knowledge, all other findings have not previously been systematically documented as explicit security vulnerabilities, although we did observe scattered fixes for related issues in some codebases. `P6` was briefly mentioned in `Tob`'s Checklist in `2023` as the corresponding entry was [added by us](https://github.com/crytic/building-secure-contracts/commit/84d54ae043773c158f1fdc8d6980bc5ccf71da91). Notably, P6 can also be viewed as a specific manifestation of the broader class of unchecked Go error handling issues; however, our work specifically highlights its severe security impact when it occurs in critical APIs such as SendCoins.


## C. PoC Summary

The [poc](./poc/) directory lists all the PoC details we conducted based on unit tests and, in some cases, we need to deploy the chain to prove our concept. For unit-test-based PoCs, the provided test code needs to be copied into the corresponding paths in the target repository before executing the tests. For certain cases, additional debugging may be required to further verify and confirm the observed differences.

We select the PoC tier that most directly and intuitively confirms each vulnerability, facilitating verification and remediation. Findings that share identical vulnerable code (e.g., fork-propagated instances) are covered by a single PoC.

<details>
<summary><b>Q: Why use Full Chain Deploy for some findings?</b></summary>

Full chain deployment is used when the vulnerability spans multiple modules (e.g., EVM + Wasm co-registration in P8) and cannot be isolated via unit test. This is the only tier that reproduces the end-to-end exploit transaction path. For single-module vulnerabilities, we prefer unit tests — using the module's `MsgServer` handler as the entry point, we construct the corresponding message to trigger the vulnerable code path directly, which is more efficient to verify and debug (e.g., setting breakpoints to observe internal state such as the gas meter before and after `callBeforeSendListener` in P1, or the unchecked error return in P6 keeper handlers).

</details>

<details>
<summary><b>Q: Why use Mock for P3 and P4 instead of UnitTest?</b></summary>

Full-chain or unit-test output (e.g., event logs) is too noisy to pinpoint the issue — dozens of events are emitted per transaction, making it difficult to observe a single missing or duplicated event. By extracting the vulnerable code pattern into a mock, we directly demonstrate the defect (missing/duplicated events). Since the surrounding code cannot compensate for these errors, this is sufficient to confirm the vulnerability.

</details>

## D. Propagation Analysis

Propagation analysis is determined through **manual code auditing** of inheritance relationships. Pairwise code similarity (`Similarity = 1 − (added + removed) / total`, via `git diff --no-index --numstat`; [source](./record/simlarity/similarity.py)) with Union-Find clustering is provided as supplementary quantitative evidence.

### Full Similarity Matrices

<details>
<summary><b>P1 Pairwise Similarity Matrix</b></summary>

| | jmes | juno | mantra | neutron | osmosis | phoneix | symphony | tower |
|---|---|---|---|---|---|---|---|---|
| **jmes** | 1.0 | 0.2222 | 0.7609 | 0.7331 | 0.7135 | 0.869 | 0.7198 | 0.7262 |
| **juno** | 0.2222 | 1.0 | 0.2456 | 0.2206 | 0.2344 | 0.2322 | 0.237 | 0.2322 |
| **mantra** | 0.7609 | 0.2456 | 1.0 | 0.7791 | 0.7156 | 0.7103 | 0.7222 | 0.7227 |
| **neutron** | 0.7331 | 0.2206 | 0.7791 | 1.0 | 0.6146 | 0.6904 | 0.6196 | 0.5863 |
| **osmosis** | 0.7135 | 0.2344 | 0.7156 | 0.6146 | 1.0 | 0.6721 | 0.981 | 0.8852 |
| **phoneix** | 0.869 | 0.2322 | 0.7103 | 0.6904 | 0.6721 | 1.0 | 0.6777 | 0.6778 |
| **symphony** | 0.7198 | 0.237 | 0.7222 | 0.6196 | 0.981 | 0.6777 | 1.0 | 0.8981 |
| **tower** | 0.7262 | 0.2322 | 0.7227 | 0.5863 | 0.8852 | 0.6778 | 0.8981 | 1.0 |

</details>

<details>
<summary><b>P2 Pairwise Similarity Matrix</b></summary>

| | jmes | juno | mantra | neutron | phoneix | symphony | tower |
|---|---|---|---|---|---|---|---|
| **jmes** | 1.0 | 0.2222 | 0.7609 | 0.7331 | 0.869 | 0.7198 | 0.7262 |
| **juno** | 0.2222 | 1.0 | 0.2456 | 0.2206 | 0.2322 | 0.237 | 0.2322 |
| **mantra** | 0.7609 | 0.2456 | 1.0 | 0.7791 | 0.7103 | 0.7222 | 0.7227 |
| **neutron** | 0.7331 | 0.2206 | 0.7791 | 1.0 | 0.6904 | 0.6196 | 0.5863 |
| **phoneix** | 0.869 | 0.2322 | 0.7103 | 0.6904 | 1.0 | 0.6777 | 0.6778 |
| **symphony** | 0.7198 | 0.237 | 0.7222 | 0.6196 | 0.6777 | 1.0 | 0.8981 |
| **tower** | 0.7262 | 0.2322 | 0.7227 | 0.5863 | 0.6778 | 0.8981 | 1.0 |

</details>

<details>
<summary><b>P3 Pairwise Similarity Matrix</b></summary>

| | jmes | juno | mantra | mitosis | neutron | osmosis | phoneix | sei | symphony | tgrade | tower |
|---|---|---|---|---|---|---|---|---|---|---|---|
| **jmes** | 1.0 | 0.2222 | 0.7609 | 0.2778 | 0.7331 | 0.7135 | 0.869 | 0.1756 | 0.2241 | 0.1357 | 0.7262 |
| **juno** | 0.2222 | 1.0 | 0.2456 | 0.2842 | 0.2206 | 0.2344 | 0.2322 | 0.25 | 0.2083 | 0.25 | 0.2322 |
| **mantra** | 0.7609 | 0.2456 | 1.0 | 0.2869 | 0.7791 | 0.7156 | 0.7103 | 0.1789 | 0.2164 | 0.1359 | 0.7227 |
| **mitosis** | 0.2778 | 0.2842 | 0.2869 | 1.0 | 0.2633 | 0.2553 | 0.2536 | 0.2483 | 0.2559 | 0.1739 | 0.2536 |
| **neutron** | 0.7331 | 0.2206 | 0.7791 | 0.2633 | 1.0 | 0.6146 | 0.6904 | 0.1538 | 0.2332 | 0.128 | 0.5863 |
| **osmosis** | 0.7135 | 0.2344 | 0.7156 | 0.2553 | 0.6146 | 1.0 | 0.6721 | 0.1447 | 0.2171 | 0.1434 | 0.8852 |
| **phoneix** | 0.869 | 0.2322 | 0.7103 | 0.2536 | 0.6904 | 0.6721 | 1.0 | 0.1572 | 0.2415 | 0.1388 | 0.6778 |
| **sei** | 0.1756 | 0.25 | 0.1789 | 0.2483 | 0.1538 | 0.1447 | 0.1572 | 1.0 | 0.144 | 0.2281 | 0.1485 |
| **symphony** | 0.2241 | 0.2083 | 0.2164 | 0.2559 | 0.2332 | 0.2171 | 0.2415 | 0.144 | 1.0 | 0.1203 | 0.2205 |
| **tgrade** | 0.1357 | 0.25 | 0.1359 | 0.1739 | 0.128 | 0.1434 | 0.1388 | 0.2281 | 0.1203 | 1.0 | 0.1388 |
| **tower** | 0.7262 | 0.2322 | 0.7227 | 0.2536 | 0.5863 | 0.8852 | 0.6778 | 0.1485 | 0.2205 | 0.1388 | 1.0 |

</details>

<details>
<summary><b>P4 Pairwise Similarity Matrix</b></summary>

| | artela | brcchain | bridgeless | cosmos-evm | cvn | de-evm | egorasmarket | guru | helios | hetu | irishub | kava | lambdavm | mezod | mtt | ollo | planq | rarimo | rollup-evm | ssc | stride | tenet | uptick | xtelabs |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| **artela** | 1.0 | 0.5921 | 0.6095 | 0.537 | 0.6225 | 0.6143 | 0.5988 | 0.607 | 0.5735 | 0.5965 | 0.2623 | 0.6002 | 0.4872 | 0.5632 | 0.5921 | 0.0479 | 0.6068 | 0.5921 | 0.5921 | 0.0818 | 0.0549 | 0.5994 | 0.111 | 0.6015 |
| **brcchain** | 0.5921 | 1.0 | 0.9421 | 0.6772 | 0.928 | 0.92 | 0.8896 | 0.9378 | 0.8193 | 0.9133 | 0.4204 | 0.9593 | 0.7117 | 0.7572 | 0.9814 | 0.0625 | 0.9642 | 0.9791 | 0.9838 | 0.0977 | 0.0621 | 0.9177 | 0.1336 | 0.961 |
| **bridgeless** | 0.6095 | 0.9421 | 1.0 | 0.7006 | 0.9635 | 0.9549 | 0.9236 | 0.9862 | 0.8519 | 0.9596 | 0.4254 | 0.964 | 0.7292 | 0.7885 | 0.9398 | 0.0622 | 0.9596 | 0.9468 | 0.9421 | 0.0974 | 0.0619 | 0.9629 | 0.1331 | 0.9329 |
| **cosmos-evm** | 0.537 | 0.6772 | 0.7006 | 1.0 | 0.7166 | 0.7014 | 0.6953 | 0.7064 | 0.7618 | 0.6885 | 0.2712 | 0.6932 | 0.5558 | 0.6485 | 0.6772 | 0.0559 | 0.6907 | 0.6818 | 0.6772 | 0.0976 | 0.0597 | 0.6872 | 0.129 | 0.6866 |
| **cvn** | 0.6225 | 0.928 | 0.9635 | 0.7166 | 1.0 | 0.9599 | 0.9301 | 0.9566 | 0.8571 | 0.9318 | 0.4274 | 0.955 | 0.7339 | 0.7649 | 0.928 | 0.0645 | 0.9482 | 0.9327 | 0.928 | 0.1006 | 0.064 | 0.9358 | 0.137 | 0.9423 |
| **de-evm** | 0.6143 | 0.92 | 0.9549 | 0.7014 | 0.9599 | 1.0 | 0.9561 | 0.9505 | 0.8552 | 0.9238 | 0.4261 | 0.9372 | 0.7278 | 0.7564 | 0.9177 | 0.0624 | 0.933 | 0.9247 | 0.92 | 0.0976 | 0.062 | 0.9392 | 0.1333 | 0.934 |
| **egorasmarket** | 0.5988 | 0.8896 | 0.9236 | 0.6953 | 0.9301 | 0.9561 | 1.0 | 0.9284 | 0.8469 | 0.9091 | 0.4128 | 0.9107 | 0.703 | 0.7455 | 0.8919 | 0.0593 | 0.9068 | 0.8941 | 0.8896 | 0.0968 | 0.0591 | 0.9288 | 0.1312 | 0.9026 |
| **guru** | 0.607 | 0.9378 | 0.9862 | 0.7064 | 0.9566 | 0.9505 | 0.9284 | 1.0 | 0.8571 | 0.9552 | 0.4227 | 0.9572 | 0.724 | 0.7873 | 0.9355 | 0.0617 | 0.9529 | 0.9424 | 0.9378 | 0.0967 | 0.0614 | 0.963 | 0.1322 | 0.9285 |
| **helios** | 0.5735 | 0.8193 | 0.8519 | 0.7618 | 0.8571 | 0.8552 | 0.8469 | 0.8571 | 1.0 | 0.8394 | 0.3748 | 0.8383 | 0.6499 | 0.6966 | 0.8193 | 0.0616 | 0.8326 | 0.8239 | 0.8193 | 0.1002 | 0.0613 | 0.8389 | 0.132 | 0.8314 |
| **hetu** | 0.5965 | 0.9133 | 0.9596 | 0.6885 | 0.9318 | 0.9238 | 0.9091 | 0.9552 | 0.8394 | 1.0 | 0.4247 | 0.935 | 0.7137 | 0.7733 | 0.9156 | 0.058 | 0.9309 | 0.9179 | 0.9133 | 0.0972 | 0.0618 | 0.9371 | 0.1329 | 0.9035 |
| **irishub** | 0.2623 | 0.4204 | 0.4254 | 0.2712 | 0.4274 | 0.4261 | 0.4128 | 0.4227 | 0.3748 | 0.4247 | 1.0 | 0.4256 | 0.4039 | 0.3577 | 0.4204 | 0.1057 | 0.4184 | 0.4204 | 0.4204 | 0.1678 | 0.121 | 0.4104 | 0.1863 | 0.4307 |
| **kava** | 0.6002 | 0.9593 | 0.964 | 0.6932 | 0.955 | 0.9372 | 0.9107 | 0.9572 | 0.8383 | 0.935 | 0.4256 | 1.0 | 0.7246 | 0.7719 | 0.9616 | 0.0629 | 0.9652 | 0.9616 | 0.9593 | 0.0983 | 0.0625 | 0.9389 | 0.1342 | 0.936 |
| **lambdavm** | 0.4872 | 0.7117 | 0.7292 | 0.5558 | 0.7339 | 0.7278 | 0.703 | 0.724 | 0.6499 | 0.7137 | 0.4039 | 0.7246 | 1.0 | 0.6282 | 0.7096 | 0.0354 | 0.7221 | 0.7075 | 0.7117 | 0.0972 | 0.0633 | 0.7078 | 0.1228 | 0.7167 |
| **mezod** | 0.5632 | 0.7572 | 0.7885 | 0.6485 | 0.7649 | 0.7564 | 0.7455 | 0.7873 | 0.6966 | 0.7733 | 0.3577 | 0.7719 | 0.6282 | 1.0 | 0.7593 | 0.0373 | 0.7733 | 0.7572 | 0.7572 | 0.0872 | 0.0512 | 0.7743 | 0.1157 | 0.7461 |
| **mtt** | 0.5921 | 0.9814 | 0.9398 | 0.6772 | 0.928 | 0.9177 | 0.8919 | 0.9355 | 0.8193 | 0.9156 | 0.4204 | 0.9616 | 0.7096 | 0.7593 | 1.0 | 0.0625 | 0.9665 | 0.9768 | 0.9814 | 0.0977 | 0.0621 | 0.92 | 0.1336 | 0.9587 |
| **ollo** | 0.0479 | 0.0625 | 0.0622 | 0.0559 | 0.0645 | 0.0624 | 0.0593 | 0.0617 | 0.0616 | 0.058 | 0.1057 | 0.0629 | 0.0354 | 0.0373 | 0.0625 | 1.0 | 0.0621 | 0.0625 | 0.0625 | 0.16 | 0.1805 | 0.0594 | 0.1106 | 0.0645 |
| **planq** | 0.6068 | 0.9642 | 0.9596 | 0.6907 | 0.9482 | 0.933 | 0.9068 | 0.9529 | 0.8326 | 0.9309 | 0.4184 | 0.9652 | 0.7221 | 0.7733 | 0.9665 | 0.0621 | 1.0 | 0.9595 | 0.9642 | 0.0972 | 0.0618 | 0.9348 | 0.1329 | 0.9412 |
| **rarimo** | 0.5921 | 0.9791 | 0.9468 | 0.6818 | 0.9327 | 0.9247 | 0.8941 | 0.9424 | 0.8239 | 0.9179 | 0.4204 | 0.9616 | 0.7075 | 0.7572 | 0.9768 | 0.0625 | 0.9595 | 1.0 | 0.9791 | 0.0977 | 0.0621 | 0.9222 | 0.1336 | 0.961 |
| **rollup-evm** | 0.5921 | 0.9838 | 0.9421 | 0.6772 | 0.928 | 0.92 | 0.8896 | 0.9378 | 0.8193 | 0.9133 | 0.4204 | 0.9593 | 0.7117 | 0.7572 | 0.9814 | 0.0625 | 0.9642 | 0.9791 | 1.0 | 0.0977 | 0.0621 | 0.9177 | 0.1336 | 0.961 |
| **ssc** | 0.0818 | 0.0977 | 0.0974 | 0.0976 | 0.1006 | 0.0976 | 0.0968 | 0.0967 | 0.1002 | 0.0972 | 0.1678 | 0.0983 | 0.0972 | 0.0872 | 0.0977 | 0.16 | 0.0972 | 0.0977 | 0.0977 | 1.0 | 0.4973 | 0.0969 | 0.2156 | 0.1006 |
| **stride** | 0.0549 | 0.0621 | 0.0619 | 0.0597 | 0.064 | 0.062 | 0.0591 | 0.0614 | 0.0613 | 0.0618 | 0.121 | 0.0625 | 0.0633 | 0.0512 | 0.0621 | 0.1805 | 0.0618 | 0.0621 | 0.0621 | 0.4973 | 1.0 | 0.0593 | 0.1349 | 0.064 |
| **tenet** | 0.5994 | 0.9177 | 0.9629 | 0.6872 | 0.9358 | 0.9392 | 0.9288 | 0.963 | 0.8389 | 0.9371 | 0.4104 | 0.9389 | 0.7078 | 0.7743 | 0.92 | 0.0594 | 0.9348 | 0.9222 | 0.9177 | 0.0969 | 0.0593 | 1.0 | 0.1314 | 0.9083 |
| **uptick** | 0.111 | 0.1336 | 0.1331 | 0.129 | 0.137 | 0.1333 | 0.1312 | 0.1322 | 0.132 | 0.1329 | 0.1863 | 0.1342 | 0.1228 | 0.1157 | 0.1336 | 0.1106 | 0.1329 | 0.1336 | 0.1336 | 0.2156 | 0.1349 | 0.1314 | 1.0 | 0.137 |
| **xtelabs** | 0.6015 | 0.961 | 0.9329 | 0.6866 | 0.9423 | 0.934 | 0.9026 | 0.9285 | 0.8314 | 0.9035 | 0.4307 | 0.936 | 0.7167 | 0.7461 | 0.9587 | 0.0645 | 0.9412 | 0.961 | 0.961 | 0.1006 | 0.064 | 0.9083 | 0.137 | 1.0 |

</details>

<details>
<summary><b>P6 Pairwise Similarity Matrix</b></summary>

| | decogit | halbornsec-ctf | haven | hero | icplaza | jackal | ollo | side-ibc | stateset | ununifi |
|---|---|---|---|---|---|---|---|---|---|---|
| **decogit** | 1.0 | 0.377 | 0.3826 | 0.3929 | 0.2628 | 0.3333 | 0.3529 | 0.3308 | 0.3673 | 0.1424 |
| **halbornsec-ctf** | 0.377 | 1.0 | 0.4425 | 0.4909 | 0.3259 | 0.3409 | 0.42 | 0.3969 | 0.4375 | 0.1807 |
| **haven** | 0.3826 | 0.4425 | 1.0 | 0.466 | 0.2656 | 0.3704 | 0.4301 | 0.3387 | 0.4494 | 0.1592 |
| **hero** | 0.3929 | 0.4909 | 0.466 | 1.0 | 0.32 | 0.3846 | 0.4667 | 0.3802 | 0.4651 | 0.1736 |
| **icplaza** | 0.2628 | 0.3259 | 0.2656 | 0.32 | 1.0 | 0.2718 | 0.3304 | 0.3288 | 0.3243 | 0.1667 |
| **jackal** | 0.3333 | 0.3409 | 0.3704 | 0.3846 | 0.2718 | 1.0 | 0.4412 | 0.303 | 0.4375 | 0.09 |
| **ollo** | 0.3529 | 0.42 | 0.4301 | 0.4667 | 0.3304 | 0.4412 | 1.0 | 0.3784 | 0.5526 | 0.1395 |
| **side-ibc** | 0.3308 | 0.3969 | 0.3387 | 0.3802 | 0.3288 | 0.303 | 0.3784 | 1.0 | 0.3738 | 0.1687 |
| **stateset** | 0.3673 | 0.4375 | 0.4494 | 0.4651 | 0.3243 | 0.4375 | 0.5526 | 0.3738 | 1.0 | 0.1279 |
| **ununifi** | 0.1424 | 0.1807 | 0.1592 | 0.1736 | 0.1667 | 0.09 | 0.1395 | 0.1687 | 0.1279 | 1.0 |

</details>

[Raw CSV](./record/simlarity/)

## E. Key Word Matching

For keyword matching, we adopt a coarse-grained, pattern-based search strategy to identify potential candidate projects. Instead of performing deep semantic analysis at this stage, we rely on direct GitHub queries with carefully designed regular expressions to quickly capture repositories that invoke specific target functions. For example:

> For P6 `bank.Sendcoinxx` project use [`/\\[\\[:space:\\]\\]\\[\\[:space:\\]\\]\\[a-zA-Z1-9\\]*\.\\[a-zA-Z1-9\\]*\.SendCoins/ language:Go`](https://github.com/search?q=%2F%5B%5B%3Aspace%3A%5D%5D%5B%5B%3Aspace%3A%5D%5D%5Ba-zA-Z1-9%5D*%5C.%5Ba-zA-Z1-9%5D*%5C.SendCoins%2F+language%3AGo&type=code&p=1).

Each PoC directory ([P1](./poc/P1/)–[P8](./poc/P8/)) contains a description of the specific matching rules and keyword patterns used for that pitfall category.

<details>
<summary><b>Q: Why is a lightweight approach sufficient?</b></summary>

Regarding detection methodology, existing static analysis tools (e.g., CodeQL and Go-analysis-based approaches) largely rely on matching predefined function signatures, which faces similar precision constraints as keyword-based matching. Furthermore, they usually operate on a predefined project, meaning the analysis results are confined to that selected codebase. Our approach instead focuses on identifying API calls that may lead to the Pitfall condition, prioritizing systematic coverage over overly restrictive exploit-specific modeling. This design aligns with our objective of identifying recurring issue classes rather than constructing a detector. By combining automated matching with manual validation, we balance scalability with practical reliability.

</details>

<details>
<summary><b>Q: Do all matches indicate vulnerabilities?</b></summary>

Crucially, we do not treat the presence of a risky pattern as a confirmed vulnerability. Some projects intentionally adopt such patterns under specific architectural constraints. For example, in `Band` (P5), although our PoC demonstrated a dirty write condition, the project's indexing mechanism prevents practical impact. Likewise, many P7 candidates deliberately restrict the message types handled by `router.handler` and relax signer checks to serve specific business logic. Among the P6 matches, from the search results, projects like `Helio` and `Injective` are matched, but they are not exploitable since the logic is in `genesis` flows or uses `SendCoinsFromModuleToAccount`, preventing attackers from directly triggering transfer failures. We excluded such contextually safe cases. Our final results therefore include only findings that we determine to remain exploitable after contextual and exploitability analysis.

The detailed exclusion records for each pitfall category in chain-registry are documented in [`./record/negative-record/`](./record/negative-record/). Since it is difficult to document the findings found directly on GitHub, we just show the positive results above.

</details>

<details>
<summary><b>Q: Could static analysis help in confirming pattern-matching results prior to manual validation?</b></summary>

We considered integrating static analysis as a second-stage confirmation step. However, general-purpose static analysis tools for Go (e.g., CodeQL, go/analysis) do not natively model the Cosmos SDK execution semantics — such as context propagation, gas meter isolation, or CacheContext event behavior. Building custom analyzers for each pitfall class would essentially duplicate the manual validation effort. Moreover, static tools typically operate on a single project at a time, whereas our GitHub-based keyword search enables cross-ecosystem discovery. We therefore chose to invest manual effort in contextual validation rather than in building bespoke static checkers.

</details>

<details>
<summary><b>Q: How conservative is the detection pipeline, and should the 65 validated cases be viewed as a lower bound?</b></summary>

Yes, the 65 validated findings represent a lower bound. Our keyword-matching step is designed for high precision at the cost of recall: we use strict patterns that target specific API calls (e.g., `.WithGasMeter`, `router.Handler`, `.CacheContext`) and manually exclude false positives. We do not claim comprehensive coverage. Projects that use alternative coding patterns, renamed functions, or custom wrappers around the vulnerable APIs may be missed. Additionally, private repositories and projects not hosted on GitHub are outside our search scope.

</details>

## F. Economically Relevant Cosmos-Based Chains Not Covered by CoinMarketCap

| Project | Mainnet / Economic Portal |
|---------|---------------------------|
| Mezo (Vaults) | [https://mezo.org/earn/vaults](https://mezo.org/earn/vaults) |
| Guru Network | [https://gurunetwork.ai/staking](https://gurunetwork.ai/staking) |
| Hetu | [https://parallel.hetu.org/p2pnetwork/earnings](https://parallel.hetu.org/p2pnetwork/earnings) |
| UnUniFi | [https://cryptorank.io/price/ununifi](https://cryptorank.io/price/ununifi) |
| Uptick | [https://www.upticknft.com/marketplace](https://www.upticknft.com/marketplace) |
| Pell Network | [https://pell.network/](https://pell.network/) |

The above projects were identified in `Table 3` as possessing economic value but were not annotated with total economic metrics due to the lack of coverage by CoinMarketCap; their mainnet portals are provided here to substantiate their economic functionality. The market capitalization of other Cosmos projects can be seen on [Mintscan](https://www.mintscan.io/).

## H. Real-world Impact Discussion and Severity Clarification

Using a conservative severity standard, we treat **P5** and **P7** as **Critical** due to direct asset-theft impact, **P2** and **P8** as **High** due to direct chain-halt impact, and **P1**, **P3**, **P4**, and **P6** as **Medium/Low** because no practical path to significant asset theft or direct chain halt has been identified. The table below lists cases that we consider relevant to the direct significant asset-impact assessment, including both cases we include and cases we exclude, together with the rationale for each decision.

| Case                                        | Scope and Rationale                                                                                                                                                                                                                                                                                                                                                                                                                           |
| ------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **P7 Secret Network**                       | **Included.** Secret’s market capitalization was approximately **$50M** when P7 was identified, and the issue had a direct asset-theft path. [SCRT](https://coinmarketcap.com/currencies/secret/)                                                                                                                                                                                                                                             |
| **P5 Evmos / Ethermint**                    | **Included.** Evmos’ market capitalization later peaked at approximately **$2B** in 2022, although the issue was identified in 2021 when market-cap records were unavailable. [EVMOS](https://coinmarketcap.com/currencies/evmos/)                                                                                                                                                                                                            |
| **P5 related EVM chains**                   | **Excluded.** These chains reused the affected EVMOS code only after P5 had already been fixed in the upstream dependency. Examples include [Cronos](https://coinmarketcap.com/currencies/cronos/), [Injective](https://coinmarketcap.com/currencies/injective/), [Kava](https://coinmarketcap.com/currencies/kava/), [ZetaChain](https://coinmarketcap.com/currencies/zetachain/), and [MANTRA](https://coinmarketcap.com/currencies/mantra/). |
| **P4 duplicated-event / fake-deposit risk** | **Excluded.** In exchanges we could register, deposits required `Tx.memos`, which blocked fake deposit exploitation through indirect transfers. Risk may still exist for systems that credit deposits from contract calls or special message handlers, where one transfer can emit multiple transfer events.                                                                                                                                                    |
| **P6 unchecked transfer failure**           | **Excluded.** The asset value is project-specific and unclear; for example, Jackal names do not have an obvious market value.                                                                                                                                                                                                                                                                                                                 |

**P1** does not directly halt chains like P2 or P8, but can cause several-fold excess gas consumption under `MaxGas`. Interestingly, when combined with another Cosmos SDK issue we identified, P1 can cause a much larger gas-consumption amplification; see the [P1 appendix](https://github.com/cosmos-bug-finder/Cosmos-Pitfalls/tree/main/poc/P1#appendix) for details.

**P3** causes event loss rather than asset theft or chain halt. This can still mislead bridges, exchanges, and indexers: critical mint or transfer records may be missed, degrading user experience and increasing reliance on centralized tracking infrastructure.




## H. Record

The [`record/`](./record/) directory contains supplementary data for this research:

- **[`acknowledgments/`](./record/acknowledgments/)** — Private disclosure acknowledgment evidence for projects that responded through private channels.
- **[`negative-record/`](./record/negative-record/)** — Exclusion records documenting projects that were matched by keyword search but excluded after contextual analysis confirmed they are not exploitable.
- **[`simlarity/`](./record/simlarity/)** — Pairwise similarity computation: the [`similarity.py`](./record/simlarity/similarity.py) script, raw CSV results, and the Go source directories used for `git diff --no-index` comparison.
- **[`chain-registry list`](https://drive.google.com/file/d/1XtHeWxdL22t-2QLRSP91nc-C4X3wUwGu/view?usp=drive_link)** — A inventory of all Cosmos chain-registr repositories, used as a part of the candidate pool for keyword matching and manual review.

---

*This research contributes to the security and robustness of the Cosmos ecosystem by identifying common pitfalls and promoting secure development practices across blockchain projects. Finally, we emphasize that our goal is not merely vulnerability detection, but the development of a structured taxonomy of customizable operations within the transaction lifecycle and the identification of concrete security pitfalls. We believe these insights can serve as a foundation for future static analysis tool development tailored specifically to the Cosmos ecosystem.*
