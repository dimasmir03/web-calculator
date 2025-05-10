package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	baseURL string
	token   string
}

func (s *APITestSuite) SetupSuite() {
	s.baseURL = "http://localhost:8080/api/v1"
	s.registerTestUser()
}

func (s *APITestSuite) registerTestUser() {
	payload := map[string]string{
		"login":    "testuser",
		"password": "testpass",
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(s.baseURL+"/register", "application/json", bytes.NewBuffer(jsonData))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	// Получаем токен
	resp, err = http.Post(s.baseURL+"/login", "application/json", bytes.NewBuffer(jsonData))
	s.Require().NoError(err)

	var result string
	json.NewDecoder(resp.Body).Decode(&result)
	s.token = result
}

func (s *APITestSuite) TestExpressionFlow() {
	// Создаем выражение
	expr := map[string]string{"expression": "2+2*2"}
	jsonData, _ := json.Marshal(expr)

	req, _ := http.NewRequest("POST", s.baseURL+"/calculate", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	// Проверяем статус
	var created struct{ ID string }
	json.NewDecoder(resp.Body).Decode(&created)

	resp, err = http.Get(s.baseURL + "/expressions/" + created.ID)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
