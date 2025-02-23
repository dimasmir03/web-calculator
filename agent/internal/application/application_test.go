package application

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dimasmir03/web-calculator-agent/pkg/workerpool"
)

func TestApplication_Run(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task" {
			json.NewEncoder(w).Encode(Task{
				Id:            "test1",
				Arg1:          10,
				Arg2:          5,
				Operation:     "Addition",
				OperationTime: 0,
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	app := &Application{
		client: NewHTTPServerClient(ts.URL),
		pool:   workerpool.NewPool(t.Context(), 2),
	}

	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
	defer cancel()

	go app.Run()
	<-ctx.Done()

}

func TestSendResult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/internal/task" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var result TaskResult
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if result.Id != "test-1" || result.Result != 15 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewHTTPServerClient(ts.URL)

	testResult := &TaskResult{
		Id:     "test1",
		Result: 15,
	}

	err := client.SendResult(context.Background(), testResult)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
