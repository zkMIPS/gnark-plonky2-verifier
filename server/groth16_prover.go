package main

import (
	"log"
	"time"
)

func proverWorkCycle(workerName string,interval uint64) {
	log.Printf("Running worker cycle")
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		log.Printf("Prover cycle started.")
	}
}