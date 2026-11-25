# Lab 6 – Producer–Consumer – Documentation

**Author:** Maksym Redchenko · **Module:** Concurrent Development, SETU Carlow, 2026/27  
**Licence:** GPL-3.0 · **Language:** Go 1.22 or newer  
**GitHub:** https://github.com/S5in/Concurrent-Development/tree/main/lab06-producer-consumer

## 1. Purpose
Producers put items into a fixed-size shared buffer and consumers take them out.
Consumers must wait when the buffer is empty, producers must wait when it is full, and only one
goroutine may change the buffer at a time.

## 2. Design
`Buffer` stores items in a slice used as a queue and uses three synchronisers:

| Synchroniser | Starts at | Meaning | Waited on by |
|---|---|---|---|
| `lock` (`sync.Mutex`) | unlocked | protects `data` | everyone, briefly |
| `items` (`Semaphore`, Lab 4) | 0 | items ready to take | consumers |
| `spaces` (`Semaphore`) | size (5) | free slots | producers |

    Put(x): spaces.Acquire → lock → append x → unlock → items.Release
    Get():  items.Acquire  → lock → take first → unlock → spaces.Release

The test uses 3 producers (10 items each, values p*100+i) and 2 consumers (15 items each).
Consumers are slower than producers so the buffer fills up and producers really have to wait.
`maxLen` records the largest number of items ever in the buffer.

## 3. Correctness
- **Mutual exclusion:** `data` is only changed while holding `lock`.
- **Empty / full:** a consumer cannot take an item until `items` > 0; a producer cannot add one
  until `spaces` > 0. The sum `items + spaces + (items in transit)` always equals the buffer size.
- **Order matters – wait on the semaphore first, then lock.** In `wrongorder/main.go`, `Get`
  locks first and then waits on `items`. With an empty buffer the consumer blocks while holding
  the lock, so no producer can get the lock to add an item → deadlock. The stack trace shows one
  consumer blocked in `Acquire` (chan receive) inside `Get` and all producers plus the other
  consumer blocked on `Mutex.Lock`. The same would happen in `Put` with a full buffer.
- **No lost or duplicated items:** each item is added once under the lock and removed once
  under the lock.

## 4. Testing and results
- `go run .`: consumed sum 3135 = expected 3135 on every run (nothing lost or duplicated);
  maximum items in the buffer = 5 = buffer size.
- `go run ./wrongorder`: `fatal error: all goroutines are asleep - deadlock!` on every run.

### Sample output (one run – order and timings change between runs)
**`go run .`**

```
consumed sum: 3135 | expected: 3135
max items in buffer: 5 | buffer size: 5
```

**`go run ./wrongorder`**

```
fatal error: all goroutines are asleep - deadlock!
```

## 5. Known issues and limitations
- A Go buffered channel `make(chan int, 5)` already is a bounded buffer with these rules built in;
  this lab builds the same behaviour from a mutex and two semaphores to show how it works.
- Consumers know in advance how many items to take; a real system would need a shutdown signal
  (for example closing a channel).
- Race detector not run (needs cgo/gcc on Windows).

## 6. Code reference
Generated from the code comments with `go doc -cmd -all -u`.

### `main.go`

```
FUNCTIONS

func main()

TYPES

type Buffer struct {
	data   []int      // the items, used as a queue (add at end, take from front)
	lock   sync.Mutex // protects data
	items  *Semaphore // number of items ready to take (consumers wait on this)
	spaces *Semaphore // number of free slots (producers wait on this)
	maxLen int        // largest number of items ever in the buffer (for checking only)
}
    Buffer is a bounded (fixed-size) buffer shared by producers and consumers.

func NewBuffer(size int) *Buffer
    NewBuffer creates an empty buffer with room for `size` items.

func (b *Buffer) Get() int
    Get removes and returns the oldest item. Waits while the buffer is empty.

func (b *Buffer) Put(x int)
    Put adds x to the buffer. Waits while the buffer is full.

type Semaphore struct {
	tokens chan struct{} // each struct{}{} in the channel is one token
}
    Semaphore is a counting semaphore built on a buffered channel (from Lab 4).
    The number of tokens in the channel is the semaphore's value.

func NewSemaphore(initial, max int) *Semaphore
    NewSemaphore creates a semaphore with `initial` tokens that can hold at most
    `max`.

func (s *Semaphore) Acquire()
    Acquire decrements the semaphore (wait / P). Blocks while the value is 0.

func (s *Semaphore) Release()
    Release increments the semaphore (signal / V), waking one waiting goroutine
    if any.
```

### `wrongorder/main.go`

```
FUNCTIONS

func main()

TYPES

type Buffer struct {
	data   []int      // the items, used as a queue (add at end, take from front)
	lock   sync.Mutex // protects data
	items  *Semaphore // number of items ready to take (consumers wait on this)
	spaces *Semaphore // number of free slots (producers wait on this)
	maxLen int        // largest number of items ever in the buffer (for checking only)
}
    Buffer is a bounded (fixed-size) buffer shared by producers and consumers.

func NewBuffer(size int) *Buffer
    NewBuffer creates an empty buffer with room for `size` items.

func (b *Buffer) Get() int
    Get removes and returns the oldest item. Waits while the buffer is empty.

func (b *Buffer) Put(x int)
    Put adds x to the buffer. Waits while the buffer is full.

type Semaphore struct {
	tokens chan struct{} // each struct{}{} in the channel is one token
}
    Semaphore is a counting semaphore built on a buffered channel (from Lab 4).
    The number of tokens in the channel is the semaphore's value.

func NewSemaphore(initial, max int) *Semaphore
    NewSemaphore creates a semaphore with `initial` tokens that can hold at most
    `max`.

func (s *Semaphore) Acquire()
    Acquire decrements the semaphore (wait / P). Blocks while the value is 0.

func (s *Semaphore) Release()
    Release increments the semaphore (signal / V), waking one waiting goroutine
    if any.
```
