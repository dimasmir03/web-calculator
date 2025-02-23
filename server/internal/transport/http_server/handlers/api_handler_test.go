package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimasmir03/web-calculator-server/pkg/calculator/cmd/calculator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetExpressionsHandler(t *testing.T) {
	e := echo.New()
	calc := calculator.NewCalculator()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := WrapperHandlerGetExpressions(calc)

	assert.NoError(t, handler(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"expressions":[]}`, rec.Body.String())

	calc.AddExpr("2+2")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	assert.NoError(t, handler(c))
	assert.Contains(t, rec.Body.String(), `"expression":"2+2"`)
}

func TestPostExpressionHandler(t *testing.T) {
	e := echo.New()
	calc := calculator.NewCalculator()

	tests := []struct {
		name       string
		expression string
		wantStatus int
	}{
		{
			name:       "valid expression",
			expression: "2+2*2",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid expression",
			expression: "2++2",
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
				fmt.Sprintf(`{"expression":"%s"}`, tt.expression)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := WrapperHandlerPostExpression(calc)
			assert.NoError(t, handler(c))
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
