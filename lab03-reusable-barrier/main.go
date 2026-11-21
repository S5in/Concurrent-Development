// Lab 3 - Reusable Barrier
// Copyright (C) 2026 Maksym Redchenko
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// --------------------------------------------
// Author: Maksym Redchenko
// Created: 21/11/2026
// Description:
//
//	Reusable (two-phase) barrier: counter + mutex + two turnstiles.
//	One turnstile is always locked, so fast goroutines cannot lap slow ones.
//	5 goroutines run 3 rounds; in every round all Part A finish before any Part B.
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

// Barrier is a REUSABLE barrier for n goroutines (two-phase / two turnstiles).
// One turnstile is always locked, so fast goroutines cannot lap slow ones.
type Barrier struct {
	n          int        // number of goroutines that must arrive
	count      int        // goroutines currently inside the barrier (protected by lock)
	lock       sync.Mutex // protects count
	turnstile1 chan bool  // door 1: buffered(1), starts EMPTY = locked
	turnstile2 chan bool  // door 2: buffered(1), starts FULL  = open
}

// NewBarrier creates a reusable barrier for n goroutines.
func NewBarrier(n int) *Barrier {
	b := &Barrier{
		n:          n,
		turnstile1: make(chan bool, 1),
		turnstile2: make(chan bool, 1),
	}
	b.turnstile2 <- true // door 2 starts open
	return b
}

// Wait blocks until all n goroutines have called Wait. Safe to call in a loop.
func (b *Barrier) Wait() {
	// ---- PHASE 1: arriving ----
	b.lock.Lock()
	b.count++
	if b.count == b.n { // last to arrive:
		<-b.turnstile2       //   lock door 2
		b.turnstile1 <- true //   open door 1
	}
	b.lock.Unlock()
	<-b.turnstile1       // pass through door 1...
	b.turnstile1 <- true // ...and let the next one through

	// ---- PHASE 2: leaving ----
	b.lock.Lock()
	b.count--
	if b.count == 0 { // last to leave:
		<-b.turnstile1       // lock door 1
		b.turnstile2 <- true // open door 2
	}
	b.lock.Unlock()
	<-b.turnstile2       // pass through door 2...
	b.turnstile2 <- true // ...and let the next one through
}

// worker runs 3 rounds; in each round it does Part A, waits at the barrier, then Part B.
func worker(id int, b *Barrier, wg *sync.WaitGroup) {
	defer wg.Done()
	for round := 1; round <= 3; round++ {
		time.Sleep(time.Duration(rand.IntN(3)) * time.Second) // arrive at random times
		fmt.Println("round", round, "Part A", id)
		b.Wait()
		fmt.Println("round", round, "Part B", id)
	}
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
