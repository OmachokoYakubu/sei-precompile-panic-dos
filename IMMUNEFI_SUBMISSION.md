# [CRITICAL] Chain-Wide DoS via Unsafe Type Assertion in Pointer Precompile

**Researcher**: Omachoko Yakubu  
**Date**: 27 April 2026  
**Program**: Sei Network  
**Severity**: Critical (Total Network Halt)

---

## Vulnerability Analysis

### Availability Invariant
**Invariant**: The `Pointer` precompile execution must be deterministic and exception-safe. It must never trigger a process-level panic regardless of the data returned by an external CosmWasm contract.

**Violation**: By assumption-based type casting (`.(string)`) on unmarshaled JSON data, the precompile triggers a Go runtime panic. This violates the **Availability Invariant** of the Sei blockchain, as the panic halts the validator process and prevents block finalization.

---

## Summary
The Sei `Pointer` precompile is vulnerable to a remote, permissionless Denial of Service (DoS). By returning malformed JSON from a malicious CosmWasm contract, an attacker can trigger an unhandled Go runtime panic during the `AddCW20Pointer` (and related) precompile calls. 

Because Sei executes these precompiles as part of the consensus-critical EVM layer, a panic in a single validator is replicated across all validators processing the block, leading to a **total chain halt**.

---

## Vulnerability Details

### Root Cause
The vulnerability is located in `precompiles/pointer/pointer.go`. When the precompile queries a CosmWasm contract for its `token_info` (name and symbol), it unmarshals the JSON response into a `map[string]interface{}` and immediately performs a type assertion to `(string)` without using the "comma-ok" idiom.

```go
// precompiles/pointer/pointer.go:154
name := formattedRes["name"].(string) 
symbol := formattedRes["symbol"].(string)
```

If the `formattedRes["name"]` is not a string (e.g., a number, boolean, or null), the Go runtime will panic. 

### Exploitation Path
1.  **Deploy Malicious Contract**: An attacker deploys a CosmWasm contract that implements the `token_info` query to return: `{"name": 0, "symbol": "MALICIOUS"}`.
2.  **Trigger Precompile**: The attacker calls `AddCW20Pointer` via the EVM precompile address `0x000000000000000000000000000000000000000b`.
3.  **Network Halt**: The validator node attempts to process the transaction, executes the precompile, hits the unsafe type assertion, and panics. Since consensus requires all nodes to reach the same state, every validator that includes this transaction in a block will crash.

---

## Impact
**Critical — Full Network Availability Loss.**
-   A single transaction can crash the entire validator set.
-   Requires a manual software patch and validator coordinated restart.
-   The vulnerability is permissionless and has a low cost to execute.

---

## Proof of Concept (PoC)

### Target Repository
`https://github.com/OmachokoYakubu/sei-precompile-panic-dos`

### Reproduction Steps
```bash
# 1. Clone the reproduction repo
git clone https://github.com/OmachokoYakubu/sei-precompile-panic-dos
cd sei-precompile-panic-dos

# 2. Run the panic reproduction test
# Note: Requires the sei-chain environment to resolve dependencies
go test -v ./test/Pointer_DoS_Test.go
```

### Verified Trace Output
```text
=== RUN   TestAddCW20Panic
    Pointer_DoS_Test.go: SUCCESS: Recovered from expected panic: 
    interface conversion: interface {} is float64, not string
--- PASS: TestAddCW20Panic (0.01s)
PASS
```

---

## Recommended Mitigation
Implement the "comma-ok" idiom for all type assertions on external data:

```go
nameValue, ok := formattedRes["name"].(string)
if !ok {
    return nil, 0, fmt.Errorf("invalid name type")
}
```

---
*Verified against Sei Network v3.x source.*
