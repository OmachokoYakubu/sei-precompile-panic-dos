# Sei Pointer Precompile: Chain-Wide DoS PoC

## Overview
This repository contains the official Proof of Concept (PoC) for a critical Denial of Service vulnerability in the Sei Network `Pointer` precompile. 

## Project Structure
- `IMMUNEFI_SUBMISSION.md`: Formal bug report for Hackerdemy.
- `test/Pointer_DoS_Test.go`: Core reproduction script (Go).
- `TRACES.txt`: Full verbose panic stack trace for forensic verification.

## Reproduction Steps

### 1. Prepare the Target Environment
Clone the official Sei repository.
```bash
git clone https://github.com/sei-protocol/sei-chain.git
cd sei-chain
```

### 2. Inject the PoC
Inject the Hackerdemy reproduction script into the local precompile directory.
```bash
# Assuming this repo is cloned adjacent to sei-chain
cp ../sei-precompile-panic-dos/test/Pointer_DoS_Test.go ./precompiles/pointer/reproduction_test.go
```

### 3. Run the Test
Execute the test suite in the target directory with verbosity.
```bash
cd precompiles/pointer/
go test -v .
```

### 4. Verify Results
A successful reproduction is confirmed when the test catches the expected Go runtime panic:
```text
=== RUN   TestAddCW20Panic
    reproduction_test.go: SUCCESS: Recovered from expected panic: 
    interface conversion: interface {} is float64, not string
--- PASS: TestAddCW20Panic (0.01s)
PASS
```

## Impact Summary
- **Severity**: Critical
- **Impact**: Full Network Halt
- **Author**: Hackerdemy

---
*Developed by Hackerdemy.*
