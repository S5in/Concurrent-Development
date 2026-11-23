// Lab 5 - Dining Philosophers (deadlock demo)
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
//	Template solution: every philosopher picks up the LEFT fork, then the RIGHT.
//	With a pause between the two forks, all 5 hold their left fork and wait
//	for their right one -> circular wait -> deadlock on every run.
//	Without the pause it usually finishes, but the deadlock is still possible.
//
// Issues:
//
//	Intentionally deadlocks (kept as evidence for the footman solution).
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

// getForks picks up the LEFT fork, then the RIGHT fork. Everyone does the same: this can deadlock.
func getForks(index int, forks []chan bool) {
	forks[index] <- true               // pick up left fork (waits if taken)
	time.Sleep(100 * time.Millisecond) // pause between forks: makes the deadlock happen every run
	forks[(index+1)%philCount] <- true // pick up right fork (waits if taken)
}

// putForks puts both forks back down.
func putForks(index int, forks []chan bool) {
	<-forks[index]
	<-forks[(index+1)%philCount]
}

// doPhilStuff is one philosopher's life: think, get forks, eat, put forks down.
func doPhilStuff(index int, wg *sync.WaitGroup, forks []chan bool) {
	defer wg.Done()
	for range meals {
		think(index)
		getForks(index, forks)
		eat(index)
		putForks(index, forks)
	}
}

func main() {
	var wg sync.WaitGroup
	forks := make([]chan bool, philCount)
	for k := range philCount {
		forks[k] = make(chan bool, 1) // one fork = channel of size 1
	}
	wg.Add(philCount)
	for n := range philCount {
		go doPhilStuff(n, &wg, forks)
	}
	wg.Wait()
	fmt.Println("everyone finished eating")
}
