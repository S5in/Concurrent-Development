# Lab 5 – Dining Philosophers – Documentation

**Author:** Maksym Redchenko · **Module:** Concurrent Development, SETU Carlow, 2026/27  
**Licence:** GPL-3.0 · **Language:** Go 1.22 or newer  
**GitHub:** https://github.com/S5in/Concurrent-Development/tree/main/lab05-dining-philosophers

## 1. Purpose
Five philosophers share five forks; each needs both neighbouring forks to eat.
Show that the template solution can deadlock and implement two solutions that cannot:
the footman and the asymmetric ("lefty") solution.

## 2. Design
Each fork is a buffered channel of size 1: picking it up is a send (blocks if the fork is
taken), putting it down is a receive. Philosopher `i` uses fork `i` (left) and fork `(i+1)%5`
(right). Each philosopher eats 3 meals (the template looped forever, so its `wg.Done()` was
never reached). All three versions contain the same 100 ms pause between picking up the first
and the second fork, which forces the worst-case timing.

| Version | Change |
|---|---|
| `deadlock` | Template: everyone picks up left then right |
| `footman` | A `Semaphore` with 4 tokens (Lab 4 type): acquire a seat before picking up forks, release after putting them down |
| `asymmetric` | The last philosopher (the "lefty") picks up right then left |

## 3. Correctness
**Requirements:** one philosopher per fork; deadlock impossible; starvation impossible;
more than one philosopher can eat at once.

- **Why the template deadlocks:** all four conditions hold – mutual exclusion, hold and wait,
  no preemption and circular wait. With the pause every philosopher holds their left fork and
  waits for their right one, which the next philosopher holds: the chain closes round the table.
  Without the pause it usually finishes, but the deadlock is still possible, so it is incorrect.
- **Footman:** at most 4 philosophers sit and there are 5 forks, so even if each holds one fork a
  fifth fork is free; at least one philosopher can get both forks, eat and release them
  (pigeonhole argument). The circular wait can never form.
- **Asymmetric:** philosopher 0 and the lefty both want fork 0 as their **first** fork. Only one
  gets it; the other waits while holding nothing and so is not part of any chain. The circle of
  "each holds one fork and waits for the next" can never close. No extra semaphore is needed and
  all 5 can sit.
- **Concurrency:** in both solutions two non-neighbours (e.g. 0 and 2) can eat at the same time.
- **Starvation:** channel waiters are woken in arrival order, so in practice every philosopher
  eats all 3 meals.

## 4. Testing and results
All versions were run repeatedly with the forced pause:
- `deadlock`: deadlocked 10/10 runs (and every run on the author's PC). Without the pause it
  finished 5/5 runs – showing that "rare" is not "impossible".
- `footman`: finished 10/10 runs.
- `asymmetric`: finished 10/10 runs.

### Sample output (one run – order and timings change between runs)
**`go run ./deadlock`**

```
Phil 1 is thinking
Phil 2 is thinking
Phil 0 is thinking
Phil 3 is thinking
Phil 4 is thinking
fatal error: all goroutines are asleep - deadlock!
```

**`go run ./asymmetric`**

```
Phil 3 is thinking
Phil 0 is thinking
Phil 1 is thinking
Phil 4 is thinking
Phil 2 is thinking
Phil 3 is eating
Phil 2 is eating
Phil 3 is thinking
Phil 1 is eating
Phil 0 is eating
Phil 2 is thinking
Phil 0 is thinking
...
```

## 5. Known issues and limitations
- Starvation-freedom relies on Go's first-come-first-served channel queues; it is not proved formally.
- Race detector not run (needs cgo/gcc on Windows).

## 6. Code reference
Generated from the code comments with `go doc -cmd -all -u`.

### `deadlock/main.go`

```
CONSTANTS

const meals = 3 // meals each philosopher eats
const philCount = 5 // philosophers (and forks)

FUNCTIONS

func doPhilStuff(index int, wg *sync.WaitGroup, forks []chan bool)
    doPhilStuff is one philosopher's life: think, get forks, eat, put forks
    down.

func eat(index int)
    eat simulates a philosopher eating for a random time.

func getForks(index int, forks []chan bool)
    getForks picks up the LEFT fork, then the RIGHT fork. Everyone does the
    same: this can deadlock.

func main()
func putForks(index int, forks []chan bool)
    putForks puts both forks back down.

func think(index int)
    think simulates a philosopher thinking for a random time.
```

### `footman/main.go`

```
CONSTANTS

const meals = 3 // meals each philosopher eats
const philCount = 5 // philosophers (and forks)

FUNCTIONS

func doPhilStuff(index int, wg *sync.WaitGroup, forks []chan bool, footman *Semaphore)
    doPhilStuff is one philosopher's life. The footman lets at most 4 sit at the
    table.

func eat(index int)
    eat simulates a philosopher eating for a random time.

func getForks(index int, forks []chan bool)
    getForks picks up the LEFT fork, then the RIGHT fork. On its own this can
    deadlock; the footman (max 4 at the table) prevents it.

func main()
func putForks(index int, forks []chan bool)
    putForks puts both forks back down.

func think(index int)
    think simulates a philosopher thinking for a random time.


TYPES

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

### `asymmetric/main.go`

```
CONSTANTS

const meals = 3 // meals each philosopher eats
const philCount = 5 // philosophers (and forks)

FUNCTIONS

func doPhilStuff(index int, wg *sync.WaitGroup, forks []chan bool)
    doPhilStuff is one philosopher's life: think, get forks, eat, put forks
    down.

func eat(index int)
    eat simulates a philosopher eating for a random time.

func getForks(index int, forks []chan bool)
    getForks: everyone picks up LEFT then RIGHT, except the last philosopher
    (the "lefty"), who picks up RIGHT then LEFT. This breaks the circular wait.

func main()
func putForks(index int, forks []chan bool)
    putForks puts both forks back down.

func think(index int)
    think simulates a philosopher thinking for a random time.
```
