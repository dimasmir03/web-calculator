package integration

import (
	"context"
	"testing"

	"github.com/dimasmir03/web-calculator-server/internal/transport/grpc/api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPC(t *testing.T) {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)
	defer conn.Close()

	client := api.NewCalculatorClient(conn)

	t.Run("Task Processing", func(t *testing.T) {
		// Тест получения задачи
		task, err := client.GetTask(context.Background(), &api.GetTaskRequest{})
		assert.NoError(t, err)

		var result float64
		switch task.Task.Operation {
		case "+":
			result = task.Task.Arg1 + task.Task.Arg2
		case "-":
			result = task.Task.Arg1 - task.Task.Arg2
		case "*":
			result = task.Task.Arg1 * task.Task.Arg2
		case "/":
			result = task.Task.Arg1 / task.Task.Arg2
		}

		// Тест отправки результата
		_, err = client.SubmitResult(context.Background(), &api.SubmitResultRequest{
			Result: result,
		})
		assert.NoError(t, err)
	})
}
