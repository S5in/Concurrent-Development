// --------------------------------------------
// Author: Maksym Redchenko
// Created: 19/11/2026
// Description:
//
//	Barrier for N goroutines: counter + mutex + turnstile (buffered channel).
//	No goroutine starts Part B until all have finished Part A.
//
// Issues:
//
//	None known.
//
// --------------------------------------------
package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// Barrier blocks goroutines until n of them have called Wait, then releases them all.
// Single-use: it cannot be reused in a loop (see Lab 3).
type Barrier struct {
	n         int        // number of goroutines that must arrive
	arrived   int        // number arrived so far (protected by lock)
	lock      sync.Mutex // protects arrived
	turnstile chan bool  // buffered(1); empty = locked
}

// NewBarrier creates a barrier for n goroutines.
// Returns a pointer so the mutex is never copied.
func NewBarrier(n int) *Barrier {
	return &Barrier{n: n, turnstile: make(chan bool, 1)}
}

// Wait blocks until all n goroutines have called Wait.
func (b *Barrier) Wait() {
	b.lock.Lock()
	b.arrived++
	if b.arrived == b.n {
		b.turnstile <- true // unlock the turnstile
	}
	b.lock.Unlock()
	<-b.turnstile       // wait for the turnstile to open
	b.turnstile <- true // let the next goroutine through
}

// worker does Part A, waits at the barrier, then does Part B.
func worker(id int, b *Barrier, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.IntN(3)) * time.Second) // arrive at random times
	fmt.Println("Part A", id)
	b.Wait() // the whole barrier is now one line
	fmt.Println("Part B", id)
}

func main() {
	var wg sync.WaitGroup
	n := 5
	b := NewBarrier(n)

	wg.Add(n)
	for id := range n {
		go worker(id, b, &wg)
	}
	wg.Wait()
}
