package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPool_addTask(t *testing.T) {
	pool := NewPool(t.Context(), 2)
	defer pool.Stop()

	var result int
	fn := func(data interface{}) error {
		val := data.(int)
		result += val
		return nil
	}

	pool.AddTask(fn, 5)
	pool.AddTask(fn, 3)

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	pool.Run()
	<-ctx.Done()

	if result != 8 {
		t.Errorf("ожидался результат 8, получен %d", result)
	}
}

func TestPool_ErrorHandling(t *testing.T) {
	pool := NewPool(t.Context(), 2)
	defer pool.Stop()

	expectedErr := errors.New("ошибка")
	fn := func(data interface{}) error {
		return expectedErr
	}

	pool.AddTask(fn, nil)

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	pool.Run()
	<-ctx.Done()
}
