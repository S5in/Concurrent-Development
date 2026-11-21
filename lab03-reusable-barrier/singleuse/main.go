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
//	Lab 2's single-use barrier used in a loop, to DEMONSTRATE the bug.
//	After round 1 the turnstile keeps its token and the counter never resets,
//	so in later rounds goroutines pass the barrier without waiting.
//
// Issues:
//
//	Intentionally broken (kept as evidence for why a reusable barrier is needed).
//
// --------------------------------------------
package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// Barrier is the SINGLE-USE barrier from Lab 2.
type Barrier struct {
	n         int        // number of goroutines that must arrive
	arrived   int        // number arrived so far (protected by lock) - never reset!
	lock      sync.Mutex // protects arrived
	turnstile chan bool  // buffered(1); empty = locked
}

// NewBarrier creates a barrier for n goroutines.
func NewBarrier(n int) *Barrier {
	return &Barrier{n: n, turnstile: make(chan bool, 1)}
}

// Wait blocks until all n goroutines have called Wait (first use only).
func (b *Barrier) Wait() {
	b.lock.Lock()
	b.arrived++
	if b.arrived == b.n {
		b.turnstile <- true // unlock the turnstile
	}
	b.lock.Unlock()
	<-b.turnstile       // wait for the turnstile to open
	b.turnstile <- true // let the next goroutine through - token is LEFT BEHIND
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
