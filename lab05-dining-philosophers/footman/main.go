// Lab 5 - Dining Philosophers (footman solution)
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
//	Footman solution: a semaphore with 4 tokens lets at most 4 of the 5
//	philosophers sit at the table. With 4 philosophers and 5 forks, at least
//	one philosopher can always get both forks (pigeonhole), so the circular
//	wait can never form -> no deadlock, even with the forced pause.
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

const philCount = 5 // philosophers (and forks)
const meals = 3     // meals each philosopher eats

// think simulates a philosopher thinking for a random time.
func think(index int) {
	time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
	fmt.Println("Phil", index, "is thinking")
}

// eat simulates a philosopher eating for a random time.
func eat(index int) {
	time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
	fmt.Println("Phil", index, "is eating")
}

// getForks picks up the LEFT fork, then the RIGHT fork.
// On its own this can deadlock; the footman (max 4 at the table) prevents it.
func getForks(index int, forks []chan bool) {
	forks[index] <- true               // pick up left fork (waits if taken)
	time.Sleep(100 * time.Millisecond) // same pause as the deadlock demo: still no deadlock
	forks[(index+1)%philCount] <- true // pick up right fork (waits if taken)
}

// putForks puts both forks back down.
func putForks(index int, forks []chan bool) {
	<-forks[index]
	<-forks[(index+1)%philCount]
}

// Semaphore is a counting semaphore built on a buffered channel (from Lab 4).
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

// doPhilStuff is one philosopher's life. The footman lets at most 4 sit at the table.
func doPhilStuff(index int, wg *sync.WaitGroup, forks []chan bool, footman *Semaphore) {
	defer wg.Done()
	for range meals {
		think(index)
		footman.Acquire() // ask the footman for a seat (may wait)
		getForks(index, forks)
		eat(index)
		putForks(index, forks)
		footman.Release() // leave the table (give the seat back)
	}
}

func main() {
	var wg sync.WaitGroup
	forks := make([]chan bool, philCount)
	for k := range philCount {
		forks[k] = make(chan bool, 1) // one fork = channel of size 1
	}
	footman := NewSemaphore(philCount-1, philCount-1) // 4 seats for 5 philosophers
	wg.Add(philCount)
	for n := range philCount {
		go doPhilStuff(n, &wg, forks, footman)
	}
	wg.Wait()
	fmt.Println("everyone finished eating")
}
