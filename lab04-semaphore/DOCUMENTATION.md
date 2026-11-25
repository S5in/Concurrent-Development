# Lab 4 – Semaphores and Multiplex – Documentation

**Author:** Maksym Redchenko · **Module:** Concurrent Development, SETU Carlow, 2026/27  
**Licence:** GPL-3.0 · **Language:** Go 1.22 or newer  
**GitHub:** https://github.com/S5in/Concurrent-Development/tree/main/lab04-semaphore

## 1. Purpose
Implement a counting semaphore type with a buffered channel (as the lecturer's
`semaphore.go` template asks: struct containing a channel, init, acquire and release) and use it
as a signal and as a multiplex that lets at most N goroutines into a critical section.

## 2. Design
**`Semaphore`.** The semaphore's value is the number of tokens (`struct{}{}`, zero bytes) in a
buffered channel.

| Operation | Lecture name | Channel operation |
|---|---|---|
| `NewSemaphore(initial, max)` | init | `make(chan struct{}, max)` + put `initial` tokens |
| `Acquire()` | wait / P / decrement | `<-tokens` (blocks while empty) |
| `Release()` | signal / V / increment | `tokens <- struct{}{}` |

The value is never read (`len` is not used). Typical uses: signal `NewSemaphore(0,1)`,
mutex `NewSemaphore(1,1)`, multiplex `NewSemaphore(N,N)`.

**`signalTest`.** A semaphore starting at 0: `main` calls `Acquire` and waits until the worker
calls `Release`.

**`multiplex`.** 12 goroutines, a semaphore with 3 tokens ("wristbands"). Each goroutine
acquires before the critical section and releases after it. An atomic counter `inside` records
how many are in the critical section at once and prints an error if it ever exceeds 3.

## 3. Correctness
- **Never more than N inside:** there are only N tokens and a goroutine needs one to enter.
- **No deadlock:** each goroutine holds at most one token and never waits for anything else while
  holding it, so there is no circular wait. Forgetting `Release` loses tokens: the first 3
  goroutines keep them forever and the other 9 block (tested → deadlock).
- **Why the value must not be read:** `if len(tokens) > 0 { Acquire() }` is check-then-act –
  another goroutine can take the token between the check and the `Acquire`.
- **Starvation:** goroutines blocked on a channel are queued in arrival order, so in practice
  every waiting goroutine eventually gets a token.

## 4. Testing and results
- Signal test: `main` always prints "got the signal" only after the worker has finished.
- Multiplex: `inside now:` never exceeds 3 and `ERROR` never appears; 12 goroutines × 0.5 s of
  work finish in about 2 s (4 batches of 3), compared with 6 s for a mutex (N = 1).
- Removing the `Release` call: `fatal error: all goroutines are asleep - deadlock!` after 3 entries.

### Sample output (one run – order and timings change between runs)
**`go run .`**

```
---- signal ----
main: waiting for worker...
worker: finished
main: got the signal
---- multiplex ----
enter 0 | inside now: 1
enter 11 | inside now: 2
enter 1 | inside now: 3
enter 2 | inside now: 3
enter 4 | inside now: 3
enter 3 | inside now: 3
enter 5 | inside now: 3
...
```

## 5. Known issues and limitations
- `max` must be at least the largest value the semaphore can reach, otherwise `Release` blocks.
- Race detector not run (needs cgo/gcc on Windows).

## 6. Code reference
Generated from the code comments with `go doc -cmd -all -u`.

### `main.go`

```
FUNCTIONS

func main()
func multiplex()
    multiplex lets at most N goroutines into the critical section at the same
    time. The semaphore starts with N tokens ("wristbands"); each goroutine
    takes one to enter and gives it back when it leaves. A mutex is the special
    case N = 1.

func signalTest()
    signalTest uses a semaphore that starts at 0 as a one-way signal: main waits
    (Acquire) until the worker signals (Release).


TYPES

type Semaphore struct {
	tokens chan struct{} // each struct{}{} in the channel is one token
}
    Semaphore is a counting semaphore built on a buffered channel. The number of
    tokens in the channel is the semaphore's value.

func NewSemaphore(initial, max int) *Semaphore
    NewSemaphore creates a semaphore with `initial` tokens that can hold at most
    `max`.

func (s *Semaphore) Acquire()
    Acquire decrements the semaphore (wait / P). Blocks while the value is 0.

func (s *Semaphore) Release()
    Release increments the semaphore (signal / V), waking one waiting goroutine
    if any.
```
