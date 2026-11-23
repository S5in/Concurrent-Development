// Lab 4 - Semaphores and Multiplex
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
// Created: 23/11/2026
// Description:
//
//	A counting semaphore type built on a buffered channel (tokens = value),
//	with Acquire (wait / P) and Release (signal / V).
//	1. signalTest: semaphore starting at 0 used as a signal between goroutines.
//	2. multiplex:  semaphore starting at N lets at most N of 12 goroutines
//	   into the critical section at once (checked with an atomic counter).
//
// Issues:
//
//	None known.
//
// --------------------------------------------

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Semaphore is a counting semaphore built on a buffered channel.
// The number of tokens in the channel is the semaphore's value.
type Semaphore struct {
	tokens chan struct{} // each struct{}{} in the channel is one token
}

// NewSemaphore creates a semaphore with `initial` tokens that can hold at most `max`.
func NewSemaphore(initial, max int) *Semaphore {
	s := &Semaphore{tokens: make(chan struct{}, max)}
	for i := 0; i < initial; i++ {
		s.tokens <- struct{}{}
	}
	return s
}

// Acquire decrements the semaphore (wait / P). Blocks while the value is 0.
func (s *Semaphore) Acquire() {
	<-s.tokens
}

// Release increments the semaphore (signal / V), waking one waiting goroutine if any.
func (s *Semaphore) Release() {
	s.tokens <- struct{}{}
}

// signalTest uses a semaphore that starts at 0 as a one-way signal:
// main waits (Acquire) until the worker signals (Release).
func signalTest() {
	// Test: use the semaphore as a SIGNAL (starts at 0)
	done := NewSemaphore(0, 1)

	go func() {
		time.Sleep(1 * time.Second) // simulate work
		fmt.Println("worker: finished")
		done.Release() // signal
	}()

	fmt.Println("main: waiting for worker...")
	done.Acquire() // wait for the signal
	fmt.Println("main: got the signal")
}

// multiplex lets at most N goroutines into the critical section at the same time.
// The semaphore starts with N tokens ("wristbands"); each goroutine takes one to
// enter and gives it back when it leaves. A mutex is the special case N = 1.
func multiplex() {
	const N = 3      // max goroutines allowed inside at once
	const total = 12 // goroutines that want to get in
	sem := NewSemaphore(N, N) // N wristbands

	var inside atomic.Int64 // how many are inside right now (measuring tool only)
	var wg sync.WaitGroup
	wg.Add(total)

	for id := range total {
		go func(id int) {
			defer wg.Done()

			sem.Acquire() // take a wristband (waits if none are left)

			// ---- critical section: at most N goroutines here ----
			now := inside.Add(1)
			fmt.Println("enter", id, "| inside now:", now)
			if now > N {
				fmt.Println("ERROR: more than", N, "inside!")
			}
			time.Sleep(500 * time.Millisecond) // simulate work
			inside.Add(-1)
			// ---- end critical section ----

			sem.Release() // give the wristband back
		}(id)
	}
	wg.Wait()
	fmt.Println("all done")
}

func main() {
	fmt.Println("---- signal ----")
	signalTest()
	fmt.Println("---- multiplex ----")
	multiplex()
}
