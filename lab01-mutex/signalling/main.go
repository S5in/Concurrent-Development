// Lab 1 - Signalling
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
//	Signalling between two goroutines using a channel.
//	Unbuffered channel: both goroutines wait for each other (rendezvous):
//	  StuffOne Part A happens before StuffTwo Part B, and
//	  StuffTwo Part A happens before StuffOne Part B.
//	With make(chan bool, 1) (buffered) it becomes a one-way signal:
//	  only StuffOne Part A before StuffTwo Part B is guaranteed.
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
	"time"
)

func main() {
	var wg sync.WaitGroup
	signal := make(chan bool) // UNBUFFERED channel: send waits for the receiver

	// doStuffOne: prints Part A, sends a signal, then prints Part B
	doStuffOne := func() {
		defer wg.Done()
		fmt.Println("StuffOne - Part A")
		signal <- true // send the signal
		fmt.Println("StuffOne - Part B")
	}

	// doStuffTwo: slow; prints Part A, waits for the signal, then prints Part B
	doStuffTwo := func() {
		defer wg.Done()
		time.Sleep(2 * time.Second) // simulate slow work
		fmt.Println("StuffTwo - Part A")
		<-signal // wait for the signal
		fmt.Println("StuffTwo - Part B")
	}

	wg.Add(2)
	go doStuffOne()
	go doStuffTwo()
	wg.Wait()
}
