// Lab 6 - Producer-Consumer (wrong order demo)
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
//	Same bounded buffer, but Get locks the mutex BEFORE waiting on `items`.
//	When the buffer is empty, a consumer waits for an item while holding the
//	lock; producers then cannot get the lock to add one -> deadlock.
//
// Issues:
//
//	Intentionally deadlocks (kept as evidence for the correct order).
//
// --------------------------------------------
package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

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

// Buffer is a bounded (fixed-size) buffer shared by producers and consumers.
type Buffer struct {
	data   []int      // the items, used as a queue (add at end, take from front)
	lock   sync.Mutex // protects data
	items  *Semaphore // number of items ready to take (consumers wait on this)
	spaces *Semaphore // number of free slots (producers wait on this)
	maxLen int        // largest number of items ever in the buffer (for checking only)
}

// NewBuffer creates an empty buffer with room for `size` items.
func NewBuffer(size int) *Buffer {
	return &Buffer{
		items:  NewSemaphore(0, size),    // empty: no items yet
		spaces: NewSemaphore(size, size), // all slots free
	}
}

// Put adds x to the buffer. Waits while the buffer is full.
func (b *Buffer) Put(x int) {
	b.spaces.Acquire() // wait for a free slot
	b.lock.Lock()
	b.data = append(b.data, x)  // add x to the end of b.data
	if len(b.data) > b.maxLen { // measuring tool: never above size
		b.maxLen = len(b.data)
	}
	b.lock.Unlock()
	b.items.Release() // announce one more item
}

// Get removes and returns the oldest item. Waits while the buffer is empty.
func (b *Buffer) Get() int {
	b.lock.Lock()       // WRONG: lock first...
	b.items.Acquire()   // ...then wait for an item WHILE HOLDING THE LOCK
	x := b.data[0]      // assign the value to x
	b.data = b.data[1:] // remove the first element
	b.lock.Unlock()
	b.spaces.Release() // announce one more free slot
	return x
}

func main() {
	const size, producers, consumers, perProducer = 5, 3, 2, 10
	total := producers * perProducer // 30 items in total
	buf := NewBuffer(size)

	var sum atomic.Int64 // sum of everything consumed (to check nothing is lost)
	var wg sync.WaitGroup
	wg.Add(producers + consumers)

	// Producers: each makes 10 numbers (producer 1 makes 100..109, etc.)
	for p := range producers {
		go func(p int) {
			defer wg.Done()
			for i := range perProducer {
				time.Sleep(time.Duration(rand.IntN(20)) * time.Millisecond) // make an item
				buf.Put(p*100 + i)
			}
		}(p)
	}

	// Consumers: each takes half of the 30 items
	for range consumers {
		go func() {
			defer wg.Done()
			for range total / consumers {
				x := buf.Get()
				sum.Add(int64(x))
				time.Sleep(time.Duration(rand.IntN(40)) * time.Millisecond) // process it (slower)
			}
		}()
	}

	wg.Wait()

	expected := 0
	for p := range producers {
		for i := range perProducer {
			expected += p*100 + i
		}
	}
	fmt.Println("consumed sum:", sum.Load(), "| expected:", expected)
	fmt.Println("max items in buffer:", buf.maxLen, "| buffer size:", size)
}
