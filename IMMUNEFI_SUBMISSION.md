# [CRITICAL] Chain-Wide DoS via Unsafe Type Assertion in Pointer Precompile

**Researcher**: Hackerdemy  
**Date**: 27 April 2026  
**Program**: Sei Network  
**Severity**: Critical (Total Network Halt)

---

## Technical Deep-Dive

### The Panic Vector
The `Pointer` precompile performs unsafe type assertions on data returned from external CosmWasm contracts. Specifically, in `pointer.go`, it assumes contract queries for `name` and `symbol` will always return string values.

### Deterministic Network Halt
Because this logic is executed during the `DeliverTx` phase of block processing, any panic is consensus-critical. If a malicious transaction is included in a block, every validator will execute the same code, hit the same panic, and exit the process. This leads to a total network halt requiring manual intervention to fix.

---

## Vulnerability Details

### Root Cause
The vulnerability is located in `precompiles/pointer/pointer.go`. When the precompile queries a CosmWasm contract for its `token_info`, it unmarshals the response into a generic map and performs a direct type assertion:

```go
// precompiles/pointer/pointer.go
name := formattedRes["name"].(string) 
symbol := formattedRes["symbol"].(string)
```

If the contract returns a JSON number or null, the Go runtime triggers a `panic`.

---

## Proof of Concept (PoC)

### Reproduction Steps

```bash
# 1. Clone the official Sei repository
git clone https://github.com/sei-protocol/sei-chain.git
cd sei-chain

# 2. Inject the Hackerdemy PoC
# Copy the provided test file from our repo into the Sei source tree
cp ../sei-precompile-panic-dos/test/Pointer_DoS_Test.go ./precompiles/pointer/reproduction_test.go

# 3. Execute the reproduction
cd precompiles/pointer/
go test -v .
```

### Expected Output
The test utilizes a `recover()` block to catch the fatal panic and confirm the vulnerability without halting the test runner. A successful reproduction logs:

```text
=== RUN   TestAddCW20Panic
    reproduction_test.go: SUCCESS: Recovered from expected panic: 
    interface conversion: interface {} is float64, not string
--- PASS: TestAddCW20Panic (0.01s)
PASS
```

---

## Impact
**Critical — Full Network Availability Loss.**
-   Bypasses the "Raised Error over Panic" hardening (PR #1507).
-   Allows permissionless, low-cost chain halts.

## Recommended Mitigation
Implement the "comma-ok" idiom for all type assertions on unmarshaled data:

```go
nameValue, ok := formattedRes["name"].(string)
if !ok {
    return nil, 0, fmt.Errorf("invalid name type")
}
```

---
*Verified by Hackerdemy against Sei Network v3.x source.*
