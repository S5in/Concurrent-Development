// Lab 2 - Rendezvous
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
//	Rendezvous of two goroutines using two one-way signals (buffered channels).
//	Each goroutine signals its arrival first, then waits for the other.
//	Guarantees: A1 happens before B2, and B1 happens before A2.
//	(Waiting first in both goroutines causes a circular wait -> deadlock.)
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

// randomSleep pauses for 0-2 seconds so the goroutines arrive in a random order.
func randomSleep() {
	time.Sleep(time.Duration(rand.IntN(3)) * time.Second)
}

func main() {
	var wg sync.WaitGroup
	aArrived := make(chan bool, 1) // signal: "A has arrived"
	bArrived := make(chan bool, 1) // signal: "B has arrived"

	wg.Add(2)

	// Goroutine A
	go func() {
		defer wg.Done()
		randomSleep()
		fmt.Println("A1")
		aArrived <- true // signal: A has arrived
		<-bArrived       // wait for B
		fmt.Println("A2")
	}()

	// Goroutine B
	go func() {
		defer wg.Done()
		randomSleep()
		fmt.Println("B1")
		bArrived <- true // signal: B has arrived
		<-aArrived       // wait for A
		fmt.Println("B2")
	}()

	wg.Wait()
}
