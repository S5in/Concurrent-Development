# Lab 3 – Reusable Barrier – Documentation

**Author:** Maksym Redchenko · **Module:** Concurrent Development, SETU Carlow, 2026/27  
**Licence:** GPL-3.0 · **Language:** Go 1.22 or newer  
**GitHub:** https://github.com/S5in/Concurrent-Development/tree/main/lab03-reusable-barrier

## 1. Purpose
Build a barrier that works every time round a loop (a reusable, two-phase barrier),
show why the Lab 2 barrier fails in a loop, and replace the mutex-protected counter with an
atomic counter.

## 2. Design
**Single-use demo (`singleuse/main.go`).** The Lab 2 barrier used for 3 rounds. After round 1
the turnstile still holds its token and `arrived` stays at n, so later rounds do not wait.

**Two-phase barrier (`main.go`).** Two turnstiles ("doors"), each a buffered channel of size 1:
`turnstile1` starts locked (empty), `turnstile2` starts open (one token). One door is always locked.

| Phase | Counter | Last goroutine does | Then everyone |
|---|---|---|---|
| 1 – arriving | `count++` | `count == n`: lock door 2, open door 1 | passes door 1 |
| 2 – leaving | `count--` | `count == 0`: lock door 1, open door 2 | passes door 2 |

After phase 2 the barrier is back in its starting state, ready for the next round.

**Atomic version (`atomic/main.go`).** `count` is an `atomic.Int64` and the mutex is removed.
`count.Add(1) == n` and `count.Add(-1) == 0` are single expressions: `Add` returns the new
value atomically, so exactly one goroutine sees n and exactly one sees 0.

## 3. Correctness
- **No lapping:** a fast goroutine that loops round reaches door 1 again, which was locked by
  the last leaver, so it must wait until all n goroutines arrive in the new round.
- **Receives inside the lock never block:** when the last goroutine arrives, every other
  goroutine has already passed door 2 and returned its token, so `<-turnstile2` returns at once.
  Likewise the last leaver finds door 1's token present (everyone decrements only after passing door 1).
- **Exactly one token** exists in total and moves door 2 → door 1 → door 2, so it is never left behind.
- **Atomic version:** splitting into `Add` then `Load` would be a check-then-act race – two
  goroutines could both see n, both try to take door 2's single token, and deadlock.

## 4. Testing and results
- Each program runs 5 goroutines for 3 rounds with random sleeps; the output was checked by a
  script that counts, for each round, whether a `Part B` line appears before all five `Part A` lines.
- Results: `singleuse` – violations in rounds 2–3 (8 in the run below); `.` and `./atomic` – 0 violations.
- Stress test of the same barrier design: 8 goroutines × 5,000 rounds × 3 runs with microsecond
  random delays – 0 violations for both the mutex and the atomic version.

### Sample output (one run – order and timings change between runs)
**`go run ./singleuse`**

```
round 1 Part A 4
round 1 Part A 1
round 1 Part A 0
round 1 Part A 2
round 1 Part A 3
round 1 Part B 4
round 2 Part A 4
round 1 Part B 1
round 2 Part A 1
round 1 Part B 0
round 1 Part B 2
round 1 Part B 3
round 2 Part B 4
round 2 Part B 1
round 3 Part A 1
round 3 Part B 1
round 2 Part A 3
round 2 Part B 3
round 3 Part A 3
round 3 Part B 3
round 2 Part A 0
round 2 Part B 0
round 2 Part A 2
round 2 Part B 2
...
```

**`go run .`**

```
round 1 Part A 4
round 1 Part A 1
round 1 Part A 3
round 1 Part A 0
round 1 Part A 2
round 1 Part B 4
round 1 Part B 1
round 2 Part A 1
round 1 Part B 3
round 2 Part A 3
round 1 Part B 0
round 1 Part B 2
...
```

## 5. Known issues and limitations
- A goroutine can start the next round's Part A while others are still printing the previous
  round's Part B; this is allowed – the barrier only orders Part A before Part B within a round.
- Race detector not run (needs cgo/gcc on Windows).

## 6. Code reference
Generated from the code comments with `go doc -cmd -all -u`.

### `singleuse/main.go`

```
FUNCTIONS

func main()
func worker(id int, b *Barrier, wg *sync.WaitGroup)
    worker runs 3 rounds; in each round it does Part A, waits at the barrier,
    then Part B.


TYPES

type Barrier struct {
	n         int        // number of goroutines that must arrive
	arrived   int        // number arrived so far (protected by lock) - never reset!
	lock      sync.Mutex // protects arrived
	turnstile chan bool  // buffered(1); empty = locked
}
    Barrier is the SINGLE-USE barrier from Lab 2.

func NewBarrier(n int) *Barrier
    NewBarrier creates a barrier for n goroutines.

func (b *Barrier) Wait()
    Wait blocks until all n goroutines have called Wait (first use only).
```

### `main.go`

```
FUNCTIONS

func main()
func worker(id int, b *Barrier, wg *sync.WaitGroup)
    worker runs 3 rounds; in each round it does Part A, waits at the barrier,
    then Part B.


TYPES

type Barrier struct {
	n          int        // number of goroutines that must arrive
	count      int        // goroutines currently inside the barrier (protected by lock)
	lock       sync.Mutex // protects count
	turnstile1 chan bool  // door 1: buffered(1), starts EMPTY = locked
	turnstile2 chan bool  // door 2: buffered(1), starts FULL  = open
}
    Barrier is a REUSABLE barrier for n goroutines (two-phase / two turnstiles).
    One turnstile is always locked, so fast goroutines cannot lap slow ones.

func NewBarrier(n int) *Barrier
    NewBarrier creates a reusable barrier for n goroutines.

func (b *Barrier) Wait()
    Wait blocks until all n goroutines have called Wait. Safe to call in a loop.
```

### `atomic/main.go`

```
FUNCTIONS

func main()
func worker(id int, b *Barrier, wg *sync.WaitGroup)
    worker runs 3 rounds; in each round it does Part A, waits at the barrier,
    then Part B.


TYPES

type Barrier struct {
	n          int64        // number of goroutines that must arrive
	count      atomic.Int64 // goroutines currently inside the barrier (no lock needed)
	turnstile1 chan bool    // door 1: buffered(1), starts EMPTY = locked
	turnstile2 chan bool    // door 2: buffered(1), starts FULL  = open
}
    Barrier is a REUSABLE barrier for n goroutines (two-phase / two turnstiles).
    One turnstile is always locked, so fast goroutines cannot lap slow ones.

func NewBarrier(n int) *Barrier
    NewBarrier creates a reusable barrier for n goroutines.

func (b *Barrier) Wait()
    Wait blocks until all n goroutines have called Wait. Safe to call in a loop.
    Add returns the new value in one indivisible step, so no mutex is needed.
```
