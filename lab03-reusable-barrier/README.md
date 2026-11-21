# Lab 3 – Reusable Barrier

**Author:** Maksym Redchenko
**Licence:** GPL-3.0 (see LICENSE in the repo root)

## How to run
Requires Go 1.22 or newer. No external libraries.

    go run ./singleuse    # Lab 2 barrier in a loop – shows the bug
    go run .              # reusable two-phase barrier (mutex)
    go run ./atomic       # reusable two-phase barrier (atomic counter)

## Files
| File | Contents |
|---|---|
| singleuse/main.go | Single-use barrier in a 3-round loop (intentionally broken) |
| main.go | Two-turnstile barrier with mutex-protected counter |
| atomic/main.go | Same barrier, counter is `atomic.Int64`, no mutex |

## Why the single-use barrier fails
After round 1 the turnstile still holds its token and `arrived` stays at n,
so in later rounds every goroutine passes straight through without waiting.

## How the two-phase barrier works
- Door 1 starts locked, door 2 starts open. One door is always locked.
- Phase 1: `count++`. The last to arrive locks door 2 and opens door 1. Everyone passes door 1.
- Phase 2: `count--`. The last to leave locks door 1 and opens door 2. Everyone passes door 2.
- The barrier is back in its start state, so it can be reused every round.

## Why the receives inside the lock are safe
When the last goroutine arrives, every other goroutine has already passed door 2 and returned
its token, so `<-turnstile2` cannot block. The same holds for the last leaver and door 1.

## Mutex vs atomic
`count.Add(1) == n` must be one expression: `Add` returns the new value atomically,
so exactly one goroutine sees n. Splitting it into `Add` then `Load` is a check-then-act race.

## Results
All three programs run 5 goroutines for 3 rounds.
- `singleuse`: correct in round 1; from round 2 goroutines print Part B without waiting.
- `.` and `./atomic`: in every round, all Part A lines appear before any Part B line.

## TO DO
- Race detector not run (needs cgo/gcc on Windows).
