package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	n := 5
	nProducer := 10
	nConsumer := 15
	product := make(chan int, n)
	var wgProducer sync.WaitGroup
	var wgConsumer sync.WaitGroup
	done := make(chan struct{})

	for i := 0; i < nProducer; i++ {
		wgProducer.Add(1)
		go func(idx int) {
			defer wgProducer.Done()
			Producer(idx, done, product)
		}(i)
	}

	for i := 0; i < nConsumer; i++ {
		wgConsumer.Add(1)
		go func(idx int) {
			defer wgConsumer.Done()
			Consumer(idx, product)
		}(i)
	}

	<-time.After(5 * time.Second)
	close(done)
	wgProducer.Wait()
	close(product)
	wgConsumer.Wait()
}

func Producer(idx int, done <-chan struct{}, products chan<- int) {
	for {
		value := rand.Int()
		select {
		case products <- value:
		case <-done:
			return
		}
		fmt.Printf("Producer %d send %d\n", idx, value)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}
}

func Consumer(idx int, product <-chan int) {
	for value := range product {
		fmt.Printf("Consumer %d receive %d\n", idx, value)
	}
}
