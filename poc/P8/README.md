# Pitfall 8

## Findings Summary

> See the [Per-Finding Index](../../README.md#per-finding-index) for the complete findings table with PoC types, fix status, and similarity values.

| # | Project | PoC Type | Issue Reference |
|---|---------|----------|-----------------|
| 63 | [Injective](https://github.com/InjectiveFoundation/injective-core) | Full Chain | [app.go#L349-L353](https://github.com/InjectiveFoundation/injective-core/blob/release/v1.18.x/injective-chain/app/app.go#L349-L353) |
| 64 | [PellNetwork](https://github.com/0xPellNetwork/aegis) | Full Chain | [app.go#L331-L341](https://github.com/0xPellNetwork/aegis/blob/main/app/app.go#L331-L341) |
| 65 | [Uptick](https://github.com/UptickNetwork/uptick) | Full Chain | [app.go#L89-L96](https://github.com/UptickNetwork/uptick/blob/main/app/app.go#L89-L96) |

---

<details>
<summary><b>Key Word Match</b></summary>

`"evmkeeper" AND "wasmd" path:app.go`

</details>

## PoC

### PoC #63–#65: Wasm+EVM (Injective, Uptick, PellNetwork)

see [Injective-case](../../high-light-findings/Injective/)
