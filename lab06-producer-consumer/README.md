# Lab 6 – Producer–Consumer

**Author:** Maksym Redchenko  
**Licence:** GPL-3.0 (see LICENSE in the repo root)  
**GitHub:** https://github.com/S5in/Concurrent-Development/tree/main/lab06-producer-consumer  
**Full repository:** https://github.com/S5in/Concurrent-Development  
**Documentation:** [DOCUMENTATION.md](DOCUMENTATION.md) (design, correctness, testing, code reference)

## How to run
Requires Go 1.22 or newer. No external libraries.

    go run .              # bounded buffer: 3 producers, 2 consumers, size 5
    go run ./wrongorder   # consumer locks before waiting – deadlocks (intentional)

## Files
| File | Contents |
|---|---|
| main.go | `Semaphore` type (Lab 4), `Buffer` with `Put`/`Get`, test with producers and consumers |
| wrongorder/main.go | Same, but `Get` locks before waiting on `items` |
| DOCUMENTATION.md | Design, correctness argument, testing/results and code reference |

## Requirements
1. Only one goroutine changes the buffer at a time (mutual exclusion).
2. A consumer waits while the buffer is empty.
3. A producer waits while the buffer is full.

## How it works
| Synchroniser | Starts at | Meaning | Waited on by |
|---|---|---|---|
| `lock` (mutex) | unlocked | protects `data` | everyone, briefly |
| `items` (semaphore) | 0 | items ready to take | consumers |
| `spaces` (semaphore) | 5 | free slots | producers |

    Put:  spaces.Acquire → lock → append → unlock → items.Release
    Get:  items.Acquire  → lock → take first → unlock → spaces.Release

Each side waits on one semaphore and releases the other.

## Why the order matters
Always **wait on the semaphore first, then lock**. In `wrongorder`, a consumer locks the
mutex and then waits on `items` while the buffer is empty. It now holds the lock while
blocked, so no producer can get the lock to add an item → deadlock.
The stack trace shows it: one consumer is blocked in `Acquire` (chan receive) inside `Get`,
while all producers and the other consumer are blocked on `Mutex.Lock`.
(The same would happen in `Put` if a producer locked before waiting on `spaces` while the buffer is full.)

## Results
- `main.go`: consumed sum 3135 = expected 3135 (nothing lost or duplicated);
  max items in buffer 5 = buffer size (producers really waited when it was full).
- `wrongorder`: `fatal error: all goroutines are asleep - deadlock!`

## Note: Go's built-in version
A Go buffered channel `make(chan int, 5)` is itself a bounded buffer: a send blocks when it
is full, a receive blocks when it is empty, and access is already mutually exclusive.
This lab builds the same behaviour from a mutex and two semaphores to show how it works.
