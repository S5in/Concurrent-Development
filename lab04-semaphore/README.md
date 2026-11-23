# Lab 4 – Semaphores and Multiplex

**Author:** Maksym Redchenko
**Licence:** GPL-3.0 (see LICENSE in the repo root)

## How to run
Requires Go 1.22 or newer. No external libraries.

    go run .    # signal test, then multiplex (12 goroutines, max 3 inside)

## Files
| File | Contents |
|---|---|
| main.go | `Semaphore` type (buffered channel), signal test, multiplex |

## The Semaphore type
- The semaphore's value is the number of tokens (`struct{}{}`) in a buffered channel.
- `NewSemaphore(initial, max)`: starts with `initial` tokens, can hold at most `max`.
- `Acquire()` (wait / P): takes a token, blocks while there are none.
- `Release()` (signal / V): puts a token back, waking one waiting goroutine.
- The value is never read (`len` is not used): reading it and then acting on it
  would be a check-then-act race, because it can change straight after the read.

| Use | Constructor |
|---|---|
| Signal | `NewSemaphore(0, 1)` |
| Mutex | `NewSemaphore(1, 1)` |
| Multiplex of N | `NewSemaphore(N, N)` |

## Multiplex
12 goroutines each take a token before the critical section and give it back after it.
An atomic counter (`inside`) records how many are inside at once; it never goes above 3.
12 goroutines × 0.5 s of work finish in about 2 s (4 batches of 3).

## Can it deadlock?
No, as long as every `Acquire` has a matching `Release`. Each goroutine holds at most one
token and never waits for anything while holding it, so there is no circular wait.
If `Release` is forgotten, the first 3 goroutines keep their tokens forever and the other 9
block in `Acquire` → `all goroutines are asleep - deadlock!` (tested).

## Starvation
Go queues goroutines blocked on a channel in arrival order, so in practice every waiting
goroutine eventually gets a token.

## TO DO
- Race detector not run (needs cgo/gcc on Windows).
