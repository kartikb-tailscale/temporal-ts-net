package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chaptersix/temporal-start-dev-ext/demo/encryption-codec-sdk-go/internal/demo"
	"github.com/chaptersix/temporal-start-dev-ext/demo/encryption-codec-sdk-go/internal/platform"
	"go.temporal.io/sdk/client"
)

func main() {
	address := flag.String("address", "127.0.0.1:7233", "Temporal frontend address")
	namespace := flag.String("namespace", "default", "Temporal namespace")
	name := flag.String("name", "temporal", "name to greet")
	codecEndpoint := flag.String("codec-endpoint", os.Getenv("TEMPORAL_CODEC_ENDPOINT"), "Remote codec endpoint")
	flag.Parse()

	c, err := platform.Dial(*address, *namespace, *codecEndpoint)
	if err != nil {
		log.Fatalf("unable to create Temporal client: %v", err)
	}
	defer c.Close()

	workflowID := fmt.Sprintf("encryption-demo-%d", time.Now().UnixNano())
	run, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: demo.TaskQueue,
	}, demo.GreetingWorkflow, *name)
	if err != nil {
		log.Fatalf("unable to execute workflow: %v", err)
	}

	var result string
	if err := run.Get(context.Background(), &result); err != nil {
		log.Fatalf("workflow failed: %v", err)
	}

	fmt.Printf("workflow_id=%s run_id=%s result=%q\n", workflowID, run.GetRunID(), result)
}
