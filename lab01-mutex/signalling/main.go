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
	signal := make(chan bool) // buffered channel used as a signal

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
