# Sei Pointer Precompile: Chain-Wide DoS PoC

## Overview
This repository contains the official Proof of Concept (PoC) for a critical Denial of Service vulnerability in the Sei Network `Pointer` precompile. The vulnerability allows an attacker to trigger a Go runtime panic across the entire validator set by returning malformed JSON data from a malicious CosmWasm contract.

## Project Structure
- `IMMUNEFI_SUBMISSION.md`: Detailed forensic report and impact analysis.
- `test/Pointer_DoS_Test.go`: Go reproduction script utilizing reflection.
- `TRACES.txt`: Captured verbose panic stack trace from a verified environment.

## Prerequisites
- Go 1.25+
- Access to the `sei-chain` source code for dependency resolution.

## Reproduction Steps

### 1. Prepare the Environment
Clone the official Sei repository and enter the precompile directory.
```bash
git clone https://github.com/sei-protocol/sei-chain.git
cd sei-chain
```

### 2. Inject the PoC
Copy the reproduction test file into the local Sei source tree.
```bash
# Assuming this repo is cloned adjacent to sei-chain
cp ../sei-precompile-panic-dos/test/Pointer_DoS_Test.go ./precompiles/pointer/reproduction_test.go
```

### 3. Execute the Exploit
Run the test suite with highest verbosity to witness the panic recovery.
```bash
cd precompiles/pointer/
go test -v .
```

### 4. Verify Results
The test will recover from the fatal interface conversion error and log the success:
```text
=== RUN   TestAddCW20Panic
    reproduction_test.go: SUCCESS: Recovered from expected panic: 
    interface conversion: interface {} is float64, not string
--- PASS: TestAddCW20Panic (0.01s)
PASS
```

## Impact Summary
- **Severity**: Critical
- **Category**: Total Network Halt (DoS)
- **Root Cause**: Unsafe type assertion `.(string)` on unmarshaled JSON data from an external contract query.

---
*Maintained by Omachoko Yakubu.*
