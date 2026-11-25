# Lab 2 – Rendezvous and Barrier – Documentation

**Author:** Maksym Redchenko · **Module:** Concurrent Development, SETU Carlow, 2026/27  
**Licence:** GPL-3.0 · **Language:** Go 1.22 or newer  
**GitHub:** https://github.com/S5in/Concurrent-Development/tree/main/lab02-barrier

## 1. Purpose
Make goroutines wait for each other at a meeting point: first two goroutines
(rendezvous), then N goroutines (barrier), packaged as a reusable-looking `Barrier` type.

## 2. Design
**Rendezvous (`rendezvous/main.go`).** Two buffered channels of size 1 act as semaphores
that start at 0: `aArrived` and `bArrived`. Each goroutine **signals its own arrival first,
then waits** for the other's signal:

    A: A1 → aArrived <- true → <-bArrived → A2
    B: B1 → bArrived <- true → <-aArrived → B2

**Barrier (`main.go`).** `Barrier` holds `n`, a counter `arrived` protected by a `sync.Mutex`,
and a `turnstile` (buffered channel of size 1, starts empty = locked).
`Wait()`:
1. Lock, `arrived++`; if `arrived == n` (last to arrive) put one token in the turnstile; unlock.
2. Take the token and immediately put it back (turnstile), so every goroutine passes one at a time.

`NewBarrier` returns `*Barrier` so the mutex inside is never copied
(the lecturer's template returned the struct by value, which `go vet` reports as a copied lock).

## 3. Correctness
- **Rendezvous:** A2 cannot run before B signalled (after B1) and B2 cannot run before A
  signalled (after A1). Because both signal before waiting, neither can block before announcing
  itself, so there is no deadlock. If both waited first, each would wait for the other forever
  (circular wait) – tested, it deadlocks.
- **Barrier:** no goroutine passes the turnstile until the n-th goroutine has incremented the
  counter, so every Part A happens before any Part B.
- **Sending inside the lock is safe:** the last arriver's send cannot block, because the
  turnstile is guaranteed empty at that moment.
- **Counter inside the lock:** `arrived++` and the `== n` check are both done under the mutex,
  otherwise updates could be lost and the last arriver might never be detected.

## 4. Testing and results
Each program was run repeatedly with random sleeps (0–2 s) before the meeting point.
- Rendezvous: in every run A1 and B1 both appear before A2 and B2.
- Barrier: in every run all five `Part A` lines appear before any `Part B` line.
- Swapping the two lines in both rendezvous goroutines (wait, then signal) produced
  `fatal error: all goroutines are asleep - deadlock!` as predicted.

### Sample output (one run – order and timings change between runs)
**`go run ./rendezvous`**

```
B1
A1
A2
B2
```

**`go run .`**

```
Part A 2
Part A 4
Part A 1
Part A 0
Part A 3
Part B 2
Part B 4
Part B 1
Part B 0
Part B 3
```

## 5. Known issues and limitations
- The barrier is **single-use**: after one use the turnstile keeps its token and the counter
  is never reset, so it does not work in a loop. This is fixed in Lab 3.
- Race detector not run (needs cgo/gcc on Windows).

## 6. Code reference
Generated from the code comments with `go doc -cmd -all -u`.

### `rendezvous/main.go`

```
FUNCTIONS

func main()
func randomSleep()
    randomSleep pauses for 0-2 seconds so the goroutines arrive in a random
    order.
```

### `main.go`

```
FUNCTIONS

func main()
func worker(id int, b *Barrier, wg *sync.WaitGroup)
    worker does Part A, waits at the barrier, then does Part B.


TYPES

type Barrier struct {
	n         int        // number of goroutines that must arrive
	arrived   int        // number arrived so far (protected by lock)
	lock      sync.Mutex // protects arrived
	turnstile chan bool  // buffered(1); empty = locked
}
    Barrier blocks goroutines until n of them have called Wait, then releases
    them all. Single-use: it cannot be reused in a loop (see Lab 3).

func NewBarrier(n int) *Barrier
    NewBarrier creates a barrier for n goroutines. Returns a pointer so the
    mutex is never copied.

func (b *Barrier) Wait()
    Wait blocks until all n goroutines have called Wait.
```
