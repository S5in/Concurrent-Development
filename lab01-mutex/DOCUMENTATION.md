# Lab 1 – Mutual Exclusion and Signalling – Documentation

**Author:** Maksym Redchenko · **Module:** Concurrent Development, SETU Carlow, 2026/27  
**Licence:** GPL-3.0 · **Language:** Go 1.22 or newer  
**GitHub:** https://github.com/S5in/Concurrent-Development/tree/main/lab01-mutex

## 1. Purpose
Show what goes wrong when several goroutines update shared data without protection,
fix it in two ways (mutex and atomic), compare their speed, and show how one goroutine can
signal another with a channel.

## 2. Design
**Shared counter (`main.go`).** 10 goroutines each add 1 to a shared `int64` 100,000 times
(expected total 1,000,000). Three versions run one after another, each with its own counter
and `sync.WaitGroup`, and each is timed with `time.Now()` / `time.Since()`:

| Version | Synchronisation | Correct? |
|---|---|---|
| `addsUnsafe` | none | No – lost updates |
| `addsMutex` | `sync.Mutex` around `*total++` | Yes |
| `addsAtomic` | `atomic.Int64.Add(1)` | Yes |

`*total++` is a read–add–write sequence. Without protection two goroutines can read the same
old value and both write back old+1, so one increment is lost (a **race condition**).
The mutex makes the increment a critical section that only one goroutine can execute at a time.
The atomic version uses a single hardware instruction, so no goroutine ever has to sleep.

**Signalling (`signalling/main.go`).** Two goroutines share an **unbuffered** channel.
StuffOne prints Part A then sends; StuffTwo sleeps, prints Part A, then receives.
An unbuffered send only completes when the receiver is ready, so both goroutines wait for
each other at that point (a two-goroutine rendezvous). With a buffered channel of size 1 the
send would not block and it would become a one-way signal.

## 3. Correctness
- **Mutex / atomic:** every increment happens inside a critical section (or as one atomic
  operation), so no update can be lost and the total is always 1,000,000.
- **Deadlock:** each goroutine holds the mutex only for one increment and never waits for
  anything else while holding it, so no circular wait is possible.
- **Signalling guarantees (unbuffered):** StuffOne Part A happens before StuffTwo Part B, and
  StuffTwo Part A happens before StuffOne Part B. With a buffered channel only the first holds.
- **Atomic vs mutex:** an atomic only protects one operation on one variable. Multi-step logic
  (check-then-act, e.g. "withdraw if balance ≥ 100") still needs a mutex.

## 4. Testing and results
Run `go run .` several times and compare the totals; run `go run ./signalling` several times
and check the guarantees above. Measured on the author's PC (Windows, 12 threads):

| Version | Total | Time |
|---|---|---|
| Unsafe | ~170,000–370,000 (changes every run) | ~3 ms |
| Mutex | 1,000,000 | ~47–51 ms |
| Atomic | 1,000,000 | ~9 ms |

The atomic version was about 5–6× faster than the mutex version; the unsafe version is fastest
but wrong.

### Sample output (one run – order and timings change between runs)
**`go run .`**

```
expected: 1000000
unsafe: got 843997 in 896.178µs
mutex : got 1000000 in 73.12738ms
atomic: got 1000000 in 16.080185ms
```

**`go run ./signalling`**

```
StuffOne - Part A
StuffTwo - Part A
StuffTwo - Part B
StuffOne - Part B
```

## 5. Known issues and limitations
- Timings vary between runs and machines (nondeterminism); the pattern, not the exact numbers, matters.
- The race detector (`go run -race`) was not run: on Windows it needs cgo and a C compiler.

## 6. Code reference
Generated from the code comments with `go doc -cmd -all -u`.

### `main.go`

```
FUNCTIONS

func addsAtomic(n int, total *atomic.Int64, wg *sync.WaitGroup)
    addsAtomic increments total n times using an atomic operation (no lock
    needed).

func addsMutex(n int, total *int64, wg *sync.WaitGroup, theLock *sync.Mutex)
    addsMutex increments total n times, protected by a mutex.

func addsUnsafe(n int, total *int64, wg *sync.WaitGroup)
    addsUnsafe increments total n times with NO protection (race condition).

func main()
```

### `signalling/main.go`

```
FUNCTIONS

func main()
```
