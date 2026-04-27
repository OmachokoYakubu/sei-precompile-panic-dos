# Chain-Wide DoS via Unsafe Type Assertion in Pointer Precompile

## Brief/Intro
The Sei `Pointer` precompile contains a critical vulnerability that allows an attacker to trigger a process-level Go runtime panic across the entire validator set. By returning malformed JSON data (a Number instead of a String) from a malicious CosmWasm contract during a precompile query, an attacker can deterministically halt the consensus mechanism, resulting in a total network-wide Denial of Service (DoS).

## Vulnerability Details
The vulnerability is located in the core logic of the `Pointer` precompile within `precompiles/pointer/pointer.go`. Specifically, when the precompile executes the `AddCW20Pointer` (or related) flow, it queries an external CosmWasm contract for its metadata (name and symbol).

The precompile unmarshals the JSON response into a generic `map[string]interface{}` and immediately performs an unsafe type assertion to `(string)` without using the "comma-ok" idiom to verify the underlying type:

```go
// precompiles/pointer/pointer.go
res, err := p.wasmdKeeper.QuerySmartSafe(ctx, cwAddr, queryBz)
if err != nil {
    return nil, 0, err
}
var formattedRes map[string]interface{}
err = json.Unmarshal(res, &formattedRes)
if err != nil {
    return nil, 0, err
}
name := formattedRes["name"].(string) // <--- CRITICAL PANIC VECTOR
symbol := formattedRes["symbol"].(string) // <--- CRITICAL PANIC VECTOR
```

In Go, `json.Unmarshal` parses JSON numbers as `float64`. If a malicious contract returns `{"name": 123}`, the type assertion `.(string)` will trigger a fatal runtime panic. Because precompiles are executed during the `DeliverTx` phase of block processing, this panic occurs synchronously across all validators, halting the chain.

## Impact Details
**Severity: Critical**
- **Impact Category**: Blockchain Halt (DoS).
- **Consequences**: Total loss of network availability. Block production is halted, preventing all transactions, liquidated-debt settlements, and bridge operations.
- **Scope Alignment**: This finding directly affects the `sei-chain` core logic and matches the "Blockchain Halt" impact in the Sei bug bounty program.
- **Exploitability**: High. Requires only permissionless CosmWasm contract deployment and a single EVM precompile call.

## References
- **Vulnerable File**: `https://github.com/sei-protocol/sei-chain/blob/main/precompiles/pointer/pointer.go`
- **Related PR (Hardening Bypass)**: PR #1507 (This vulnerability bypasses the "Raised Error over Panic" hardening implemented in previous versions).

## Proof of Concept

### Reproduction Steps

```bash
# 1. Clone the Hackerdemy reproduction repository
git clone https://github.com/OmachokoYakubu/sei-precompile-panic-dos.git
cd sei-precompile-panic-dos

# 2. Clone the official Sei repository (adjacent to PoC)
cd ..
git clone https://github.com/sei-protocol/sei-chain.git
cd sei-chain

# 3. Inject the Hackerdemy PoC
# Copy the provided test file from our repo into the Sei source tree
cp ../sei-precompile-panic-dos/test/Pointer_DoS_Test.go ./precompiles/pointer/reproduction_test.go

# 4. Execute the reproduction
cd precompiles/pointer/
go test -v .
```

### Expected Output
The test utilizes a `recover()` block to catch the fatal panic and confirm the vulnerability. A successful reproduction logs:

```text
=== RUN   TestAddCW20Panic
    reproduction_test.go: SUCCESS: Recovered from expected panic: 
    interface conversion: interface {} is float64, not string
--- PASS: TestAddCW20Panic (0.01s)
PASS
```

---
*Submitted by Hackerdemy.*
