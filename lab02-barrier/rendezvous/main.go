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
		aArrived <- true
		<-bArrived
		fmt.Println("A2")
	}()

	// Goroutine B
	go func() {
		defer wg.Done()
		randomSleep()
		fmt.Println("B1")
		bArrived <- true
		<-aArrived
		fmt.Println("B2")
	}()

	wg.Wait()
}
