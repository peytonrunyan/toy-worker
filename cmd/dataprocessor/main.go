package main

import (
	c "dataprocessor/internal/common"
	r "dataprocessor/internal/repository"
	w "dataprocessor/internal/worker"
	"flag"
	"log"
	"sync"
)

func main() {
	setProducer := flag.Bool("produce", false, "Whether you want the worker to produce (set flag) or consume (default value when flag not set)")
	flag.Parse()

	rmq, err := r.NewRMQ(c.RMQ_CONN_STRING)
	if err != nil {
		log.Fatal(err)
	}
	defer rmq.CleanUp()

	worker := w.NewRabbitWorker(rmq)
	var wg sync.WaitGroup

	wg.Add(1)

	switch {
	case *setProducer:
		worker.Produce(&wg)
	default:
		worker.Consume(&wg)
	}

	wg.Wait()
}
