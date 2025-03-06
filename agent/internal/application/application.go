package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dimasmir03/web-calculator-agent/pkg/workerpool"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerURL      string
	ComputingPower int
	RetryDelay     time.Duration
}

type Task struct {
	Id            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type TaskResult struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
}

type ServerClient interface {
	GetTask(ctx context.Context) (*Task, error)
	SendResult(ctx context.Context, result *TaskResult) error
}

type Application struct {
	config   *Config
	client   *HTTPServerClient
	pool     *workerpool.Pool
	stopChan chan os.Signal
}

func NewApplication() (*Application, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("ошибка загрузки .env")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфига: %w", err)
	}

	client := NewHTTPServerClient(cfg.ServerURL)

	ctx := context.Background()
	pool := workerpool.NewPool(ctx, cfg.ComputingPower)

	return &Application{
		config:   cfg,
		client:   client,
		pool:     pool,
		stopChan: make(chan os.Signal, 1),
	}, nil
}

func (a *Application) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal.Notify(a.stopChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	a.pool.Run()
	defer a.pool.Stop()

	go a.taskProcessor(ctx)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case sig := <-a.stopChan:
		slog.Info("прекрасный выключение:", "signal", sig)
		return nil
	}
}

func (a *Application) taskProcessor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			task, err := a.client.GetTask(ctx)
			if err != nil {
				a.handleTaskError(err)
				continue
			}

			a.processTask(ctx, task)
		}
	}
}

func (a *Application) processTask(ctx context.Context, task *Task) {
	slog.Info("обработка таски", "id", task.Id, slog.Any("task", task))

	processTask := workerpool.TaskFunc(func(data interface{}) error {
		t := data.(*Task)
		result, err := a.calculate(t)
		if err != nil {
			return err
		}
		return a.client.SendResultWithRetry(ctx, result)
	})

	a.pool.AddTask(processTask, task)
}

func (a *Application) calculate(task *Task) (*TaskResult, error) {
	startTime := time.Now()
	slog.Info("обработка началась", slog.Any("task", task))
	res := new(TaskResult)

	var result float64
	switch task.Operation {
	case "Addition":
		result = task.Arg1 + task.Arg2
	case "Substraction":
		result = task.Arg1 - task.Arg2
	case "Multiplication":
		result = task.Arg1 * task.Arg2
	case "Division":
		if task.Arg2 == 0 {
			return nil, errors.New("деление на 0")
		}
		result = task.Arg1 / task.Arg2
	default:
		return nil, fmt.Errorf("неизвестный операнд: %s", task.Operation)
	}
	res.Id = task.Id
	res.Result = result
	slog.Info("таска обработана",
		"id", task.Id,
		"time", time.Since(startTime),
		slog.Any("task", task),
	)

	return res, nil
}

func (a *Application) handleTaskError(err error) {
	slog.Error("ошибка получения таски", "error", err)
	time.Sleep(a.config.RetryDelay)
}

func LoadConfig() (*Config, error) {
	retryDelay, err := strconv.Atoi(os.Getenv("TIME_WAIT_MS"))
	if err != nil {
		retryDelay = 1000
	}

	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		return nil, errors.New("ошибка загрузки env COMPUTING_POWER")
	}

	return &Config{
		ServerURL:      os.Getenv("SERVER_URL"),
		ComputingPower: computingPower,
		RetryDelay:     time.Duration(retryDelay) * time.Millisecond,
	}, nil
}

type HTTPServerClient struct {
	baseURL    string
	client     *http.Client
	retryDelay time.Duration
}

func NewHTTPServerClient(baseURL string) *HTTPServerClient {
	return &HTTPServerClient{
		baseURL:    baseURL,
		client:     &http.Client{Timeout: 5 * time.Second},
		retryDelay: 1 * time.Second,
	}
}

func (s *HTTPServerClient) GetTask(ctx context.Context) (*Task, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/internal/task", nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("что-то пошло не так: code=%d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *HTTPServerClient) SendResult(ctx context.Context, result *TaskResult) error {
	jsondata, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("ошибка кодирования в json: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL+"/internal/task",
		bytes.NewBuffer(jsondata),
	)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	return nil
}

func (s *HTTPServerClient) SendResultWithRetry(ctx context.Context, result *TaskResult) error {
	err := s.SendResult(ctx, result)
	if err == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(s.retryDelay):
	}
	return fmt.Errorf("max retries exceeded")
}
