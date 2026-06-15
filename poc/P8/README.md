# Pitfall 8

## Findings Summary

> See the [Per-Finding Index](../../README.md#per-finding-index) for the complete findings table with PoC types, fix status, and similarity values.

| # | Project | PoC Type | Issue Reference |
|---|---------|----------|-----------------|
| 63 | [Injective](https://github.com/InjectiveFoundation/injective-core) | Full Chain | [app.go#L349-L353](https://github.com/InjectiveFoundation/injective-core/blob/release/v1.18.x/injective-chain/app/app.go#L349-L353) |
| 64 | [Uptick](https://github.com/UptickNetwork/uptick) | Full Chain | [app.go#L89-L96](https://github.com/UptickNetwork/uptick/blob/main/app/app.go#L89-L96) |
| 65 | [Althea](https://github.com/AltheaFoundation/althea-L1) | Full Chain | [-](https://github.com/AltheaFoundation/althea-L1/commit/ced08fcbdb5c0c1f2b16fa79c62cd0301c64ad91) |

---

<details>
<summary><b>Key Word Match</b></summary>

`"evmkeeper" AND "wasmd" path:app.go`

</details>

## PoC

### PoC #63–#65: Wasm+EVM (Injective, Uptick, PellNetwork)

see [Injective-case](../../high-light-findings/Injective/)
