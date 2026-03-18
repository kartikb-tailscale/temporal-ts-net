package demo

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "encryption-demo-task-queue"

func GreetingWorkflow(ctx workflow.Context, name string) (string, error) {
	opts := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, opts)

	var result string
	if err := workflow.ExecuteActivity(ctx, ComposeGreetingActivity, name).Get(ctx, &result); err != nil {
		return "", err
	}
	return result, nil
}

func ComposeGreetingActivity(_ context.Context, name string) (string, error) {
	return fmt.Sprintf("hello %s", name), nil
}
