package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		result := job + 1

		// Send the result to the receiver channel
		results <- result

		fmt.Printf("Worker %d processed job %d, result: %d\n", id, job, result)
	}
}

func main() {
	numWorkers := 3
	numJobs := 10

	jobs := make(chan int)
	results := make(chan int)

	// Created a wait group to ensure all workers finish
	var wg sync.WaitGroup

	// Created the workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i+1, jobs, results, &wg)
	}

	// Send jobs to the workers
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Waited for all workers to finish in go routine as to not block the main thread from receiving results and closing the results channel as if not done the main thread loop will infinitely wait for results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Receive and print the results from the workers
	for result := range results {
		fmt.Printf("Received result: %d\n", result)
	}
}
