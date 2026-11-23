# Lab 5 – Dining Philosophers

**Author:** Maksym Redchenko
**Licence:** GPL-3.0 (see LICENSE in the repo root)

## How to run
Requires Go 1.22 or newer. No external libraries.

    go run ./deadlock      # template behaviour – deadlocks every run (intentional)
    go run ./footman       # footman solution (max 4 at the table)
    go run ./asymmetric    # asymmetric solution (one "lefty")

All three versions use the same forced 100 ms pause between picking up the first
and second fork. This makes the deadlock happen every time in `deadlock`, so the
two solutions are tested under the worst timing.

## Files
| File | Contents |
|---|---|
| deadlock/main.go | Everyone picks up left then right fork → circular wait → deadlock |
| footman/main.go | Semaphore with 4 tokens (Lab 4 type) lets at most 4 sit at the table |
| asymmetric/main.go | Last philosopher picks up right then left fork |

## Requirements
1. Only one philosopher can hold a fork at a time (each fork is a channel of size 1).
2. Deadlock must be impossible.
3. Starvation must be impossible.
4. More than one philosopher must be able to eat at the same time.

## Why the template deadlocks
Every philosopher holds their left fork and waits for their right fork, which the next
philosopher holds. The chain closes round the table: a **circular wait**.
All four deadlock conditions hold: mutual exclusion, hold and wait, no preemption, circular wait.
Without the pause it usually finishes, but the deadlock is still possible, so it is wrong.

## Footman solution – why it cannot deadlock
At most 4 philosophers sit at the table, and there are 5 forks. Even if each of the 4
holds one fork, one fork is free, so at least one philosopher can get both forks, eat,
and put them down. The circular wait can never form. Two non-neighbours can still eat
at the same time.

## Asymmetric solution – why it cannot deadlock
Philosopher 0 and the lefty (philosopher 4) both want fork 0 as their first fork.
Only one of them gets it; the other waits while holding nothing, so it is not part of
any chain. The circle of "each holds one fork and waits for the next" can never close.
No extra semaphore is needed and all 5 can sit at the table.

## Starvation
Goroutines blocked on a channel are woken in arrival order, so a waiting philosopher
eventually gets the fork (or seat). In practice no philosopher starves; each eats all 3 meals.

## Results
- deadlock: `fatal error: all goroutines are asleep - deadlock!` on every run.
- footman: finished on every run (tested 10/10).
- asymmetric: finished on every run (tested 10/10).

## TO DO
- Race detector not run (needs cgo/gcc on Windows).
