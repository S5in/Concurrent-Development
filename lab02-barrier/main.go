// Lab 2 - Barrier
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
// Created: 19/11/2026
// Description:
//
//	Barrier for N goroutines: counter + mutex + turnstile (buffered channel).
//	No goroutine starts Part B until all have finished Part A.
//
// Issues:
//
//	Single-use only: not reusable in a loop (fixed in Lab 3).
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
