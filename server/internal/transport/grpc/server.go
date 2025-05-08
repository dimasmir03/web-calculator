package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/dimasmir03/web-calculator-server/internal/calculator/cmd/calculator"
	"github.com/dimasmir03/web-calculator-server/internal/transport/grpc/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	api.UnimplementedCalculatorServer
	mu      sync.Mutex
	Calc    *calculator.Calculator
	tasks   map[string]*api.Task
	results chan *api.SubmitResultRequest
}

func (s *Server) GetTask(
	ctx context.Context,
	req *api.GetTaskRequest,
) (
	*api.GetTaskResponse,
	error,
) {
	expr := s.Calc.GetSimpleExpr()
	fmt.Println(expr)
	if string(expr.Id) == "" {
		return nil, status.Error(codes.NotFound, "нету задач")
	}
	var tt int
	switch expr.Op {
	case "+":
		tt, _ = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	case "-":
		tt, _ = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	case "*":
		tt, _ = strconv.Atoi(os.Getenv("TIME_MULTIPLICATION_MS"))
	case "/":
		tt, _ = strconv.Atoi(os.Getenv("TIME_DIVISION_MS"))
	}
	task := &api.Task{
		Id:            string(expr.Id),
		Arg1:          expr.A.(float64),
		Arg2:          expr.B.(float64),
		Operation:     expr.Op,
		OperationTime: int64(tt),
	}
	fmt.Println("CHECK", task)
	return &api.GetTaskResponse{Task: task}, nil
}

func (s *Server) SubmitResult(
	ctx context.Context,
	req *api.SubmitResultRequest,
) (
	*api.SubmitResultResponse,
	error,
) {
	slog.Debug("Результат: ", req)
	if err := s.Calc.SetSimpleExprResult(req.Id, req.Result, ""); err != nil {
		return &api.SubmitResultResponse{Success: false}, status.Error(codes.NotFound, err.Error())
	}
	return &api.SubmitResultResponse{Success: true}, nil
}
