# Lab 1 – Mutual Exclusion and Signalling

**Author:** Maksym Redchenko
**Licence:** GPL-3.0 (see LICENSE in the repo root)

## How to run
Requires Go 1.22 or newer. No external libraries.
    go run .               # race demo + mutex + atomic, with timings
    go run ./signalling    # signalling between two goroutines

## Files
| File | Contents |
|---|---|
| main.go | Shared counter: unsafe (race), mutex and atomic versions, timed |
| signalling/main.go | One goroutine signals another with a channel |

## Results (my PC)
| Version | Total | Time |
|---|---|---|
| Unsafe | ~170,000–300,000 (varies, wrong) | ~3 ms |
| Mutex | 1,000,000 | ~47–51 ms |
| Atomic | 1,000,000 | ~9 ms |

## What I learned
- `total++` is read-add-write; goroutines can interleave, so updates are lost (race condition).
- A mutex makes the critical section run one goroutine at a time: correct but slower (contention).
- Atomics are faster but only protect a single operation on a single variable;
  multi-step logic (check-then-act, transfers) still needs a mutex.
- Unbuffered channel = both sides wait (rendezvous); buffered = one-way signal.
