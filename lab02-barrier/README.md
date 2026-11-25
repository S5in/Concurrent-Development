# Lab 2 – Rendezvous and Barrier

**Author:** Maksym Redchenko  
**Licence:** GPL-3.0 (see LICENSE in the repo root)  
**GitHub:** https://github.com/S5in/Concurrent-Development/tree/main/lab02-barrier  
**Full repository:** https://github.com/S5in/Concurrent-Development  
**Documentation:** [DOCUMENTATION.md](DOCUMENTATION.md) (design, correctness, testing, code reference)

## How to run
Requires Go 1.22 or newer. No external libraries.

    go run ./rendezvous    # rendezvous of two goroutines
    go run .               # barrier for 5 goroutines

## Files
| File | Contents |
|---|---|
| rendezvous/main.go | Two goroutines wait for each other (two one-way signals) |
| main.go | Barrier type: counter + mutex + turnstile, used by 5 goroutines |
| DOCUMENTATION.md | Design, correctness argument, testing/results and code reference |

## How the barrier works
1. Each goroutine locks the mutex, increments `arrived`, and the last one (arrived == n)
   puts a token in the turnstile. The send cannot block: the turnstile is empty at that point.
2. Every goroutine then takes the token and puts it back (turnstile), so all pass one by one.

## What I learned
- Rendezvous: signal first, then wait. Waiting first in both goroutines → circular wait → deadlock.
- Never block while holding a lock: the other goroutines can't get the lock, so everything freezes.
- Return *Barrier (pointer) so the mutex is never copied.

