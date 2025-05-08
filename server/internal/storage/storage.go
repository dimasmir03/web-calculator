package storage

//
// import (
// 	"sync"
//
// 	"github.com/dimasmir03/web-calculator-server/internal/calculator/cmd/calculator"
// )
//
// type Store struct {
// 	calc *calculator.Calculator
// 	m    sync.Mutex
// }
//
// func NewStore() *Store {
// 	return &Store{
// 		calc: calculator.NewCalculator(),
// 		m:    sync.Mutex{},
// 	}
// }
//
// func (s *Store) addExpression(expression string) (string, error) {
// 	s.m.Lock()
// 	defer s.m.Unlock()
// 	id, err := s.calc.AddExpr(expression)
// 	if err != nil {
// 		return "", err
// 	}
// 	return id, nil
// }
//
// func (s *Store) getExpressions() []calculator.Expr {
// 	s.m.Lock()
// 	defer s.m.Unlock()
// 	return s.calc.GetExpressionsStatus()
// }
