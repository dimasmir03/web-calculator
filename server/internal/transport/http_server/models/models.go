package models

type (
	// Error response
	// @name ErrorResponse
	ErrResponse struct {
		Message string `json:"message"`
	}

	// Calculate request
	// @name CalculateRequest
	CalculateRequest struct {
		Expression string `json:"expression" example:"2+2*2"`
	}

	// Calculate response
	// @name CalculateResponse
	CalculateResponse struct {
		Id string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	}

	// Expression item
	// @name Expression
	Expression struct {
		Id         string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
		Expression string `json:"expression" example:"2+2*2"`
		Status     string `json:"status" example:"completed"`
		Result     string `json:"result" example:"6"`
	}

	// Expressions list
	// @name ExpressionsResponse
	ExpressionsResponse struct {
		Expressions []Expression `json:"expressions"`
	}

	// ExpressionResponse GET 'localhost/api/v1/expressions/:id'
	ExpressionResponse struct {
		Expressions Expression `json:"expression"`
	}

	// Task GET 'localhost/internal/task'

	TaskResponse struct {
		Id            string  `json:"id"`
		Arg1          float64 `json:"arg1"`
		Arg2          float64 `json:"arg2"`
		Operation     string  `json:"operation"`
		OperationTime int     `json:"operation_time"`
	}

	// TaskResult POST 'localhost/internal/task'

	TaskResultRequest struct {
		Id     string  `json:"id"`
		Result float64 `json:"result,omitempty"`
		Error  string  `json:"error,omitempty"`
	}
)
