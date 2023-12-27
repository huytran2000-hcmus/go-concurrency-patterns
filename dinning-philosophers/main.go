package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var forks [5]sync.Mutex
	go Philosopher(0, &forks[0], &forks[4])
	go Philosopher(1, &forks[0], &forks[1])
	go Philosopher(2, &forks[1], &forks[2])
	go Philosopher(3, &forks[2], &forks[3])
	go Philosopher(4, &forks[3], &forks[4])

	select {}
}

func Philosopher(index int, firstFork, secondFork *sync.Mutex) {
	for {

		fmt.Printf("Philosopher %d is thinking\n", index)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		firstFork.Lock()
		secondFork.Lock()
		fmt.Printf("Philosopher %d is eating\n", index)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		secondFork.Unlock()
		firstFork.Unlock()
	}
}
