package main

import (
	"flag"
	"log"
	"os"

	"github.com/chaptersix/temporal-start-dev-ext/demo/encryption-codec-sdk-go/internal/demo"
	"github.com/chaptersix/temporal-start-dev-ext/demo/encryption-codec-sdk-go/internal/platform"
	"go.temporal.io/sdk/worker"
)

func main() {
	address := flag.String("address", "127.0.0.1:7233", "Temporal frontend address")
	namespace := flag.String("namespace", "default", "Temporal namespace")
	codecEndpoint := flag.String("codec-endpoint", os.Getenv("TEMPORAL_CODEC_ENDPOINT"), "Remote codec endpoint")
	flag.Parse()

	c, err := platform.Dial(*address, *namespace, *codecEndpoint)
	if err != nil {
		log.Fatalf("unable to create Temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, demo.TaskQueue, worker.Options{})
	w.RegisterWorkflow(demo.GreetingWorkflow)
	w.RegisterActivity(demo.ComposeGreetingActivity)

	log.Printf("worker started on task queue %q", demo.TaskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("worker exited with error: %v", err)
	}
}
