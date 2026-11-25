// Lab 1 - Mutual Exclusion
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

//--------------------------------------------
// Author: Maksym Redchenko
// Created: 17/11/2026
// Description:
//   10 goroutines each add 1 to a shared counter 100,000 times.
//   Version 1 (unsafe) shows the race condition: updates are lost.
//   Version 2 (mutex) and Version 3 (atomic) fix it; both are timed.
// Issues:
//   None known.
//--------------------------------------------

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// addsUnsafe increments total n times with NO protection (race condition).
func addsUnsafe(n int, total *int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		*total++ // read-add-write can interleave between goroutines
	}
}

// addsMutex increments total n times, protected by a mutex.
func addsMutex(n int, total *int64, wg *sync.WaitGroup, theLock *sync.Mutex) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		theLock.Lock() // enter critical section
		*total++
		theLock.Unlock() // leave critical section
	}
}

// addsAtomic increments total n times using an atomic operation (no lock needed).
func addsAtomic(n int, total *atomic.Int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		total.Add(1) // read-add-write as one indivisible step
	}
}

func main() {
	routines, n := 10, 100000
	fmt.Println("expected:", routines*n)

	// ---- Version 1: unsafe (shows the race) ----
	var wg1 sync.WaitGroup
	var unsafeTotal int64
	start := time.Now()
	wg1.Add(routines)
	for range routines {
		go addsUnsafe(n, &unsafeTotal, &wg1)
	}
	wg1.Wait()
	fmt.Println("unsafe: got", unsafeTotal, "in", time.Since(start))

	// ---- Version 2: mutex ----
	var wg2 sync.WaitGroup
	var theLock sync.Mutex
	var mutexTotal int64
	start = time.Now()
	wg2.Add(routines)
	for range routines {
		go addsMutex(n, &mutexTotal, &wg2, &theLock)
	}
	wg2.Wait()
	fmt.Println("mutex : got", mutexTotal, "in", time.Since(start))

	// ---- Version 3: atomic ----
	var wg3 sync.WaitGroup
	var atomicTotal atomic.Int64
	start = time.Now()
	wg3.Add(routines)
	for range routines {
		go addsAtomic(n, &atomicTotal, &wg3)
	}
	wg3.Wait()
	fmt.Println("atomic: got", atomicTotal.Load(), "in", time.Since(start))
}
