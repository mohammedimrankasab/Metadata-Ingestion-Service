package main

import (
	"time"

	"github.com/mohammedimrankasab/go-concurrency-examples/internal/logger"
)

func producer(ch chan<- string) {
	workspaces := []string{
		"Finance",
		"Sales",
		"Marketing",
		"HR",
		"Engineering",
	}

	for _, workspace := range workspaces {

		ch <- workspace
	}
	close(ch)

}

func consumer(ch <-chan string) {

	for workspace := range ch {
		// Implementation for ingesting workspace
		logger.Info("Ingesting workspace: " + workspace)
		time.Sleep(2 * time.Second)
		logger.Info("Finished ingesting workspace: " + workspace)
	}
}

func main() {

	ch := make(chan string)

	go producer(ch)
	consumer(ch)
}
