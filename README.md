# Sei Pointer Precompile DoS — Reproduction Guide

## Overview
This repository contains the proof of concept for a critical Denial of Service (DoS) vulnerability in the Sei `Pointer` precompile. The PoC demonstrates a Go runtime panic triggered by unsafe type assertions on malformed JSON data returned from a contract query.
This repository contains a proof of concept (PoC) for a critical Denial of Service (DoS) vulnerability found in the Sei `Pointer` precompile. It demonstrates how malformed JSON data returned from a contract query can trigger an unhandled Go runtime panic.

## Prerequisites
- Go 1.25+
- `seid` source code (`sei-chain` repo)

## Project Structure
- `IMMUNEFI_SUBMISSION.md`: Formal bug report and technical deep-dive.
- `test/Pointer_DoS_Test.go`: The Go-based panic reproduction.
- `TRACES.txt`: Captured panic stack trace from a verified run.

## Reproduction Steps

### 1. Clone the Target Repository
```bash
git clone https://github.com/sei-protocol/sei-chain.git
cd sei-chain
```

```

### 4. Verify the Panic
The test will output a success message upon recovering from the expected runtime panic:
```
=== RUN   TestAddCW20PanicReflection
    exploit_test.go:78: SUCCESS: Recovered from expected panic: interface conversion: interface {} is float64, not string
--- PASS: TestAddCW20PanicReflection (0.01s)
```

## Vulnerability Analysis (Beanstalk Standard)
**Invariant**: `AddCW20Pointer(addr)` must never result in a process-level panic regardless of the data returned by `addr`.
**Violation**: By returning a JSON number where a string is expected, the precompile triggers a Go `panic`, violating the availability invariant of the blockchain node.
