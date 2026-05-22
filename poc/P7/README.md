# Pitfall 7

## Findings Summary

> See the [Per-Finding Index](../../README.md#per-finding-index) for the complete findings table with PoC types, fix status, and similarity values.

| # | Project | PoC Type | Issue Reference |
|---|---------|----------|-----------------|
| 62 | [Secret Network](https://github.com/scrtlabs/SecretNetwork) | Full Chain | Screenshot evidence |

---

<details>
<summary><b>Key Word Match</b></summary>

`/cosmos\/cosmos-sdk/ AND /\.router\.Handler\(/ NOT /GetMsgV1Signers/ NOT /GetSigners/`

Here we manually excluded cases that could only send fixed message types because they were difficult to utilize or were issued by the governance module without signature verification.

</details>

## PoC

### PoC #62: Secret Network

see [secret-case](../../high-light-findings/secret/)