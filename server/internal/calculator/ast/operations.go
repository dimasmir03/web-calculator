package ast

var (
	operationsStr = []string{"Invalid", "+", "-", "*", "/"}
)

type Operation uint8

const (
	Invalid Operation = iota

	Addition
	Substraction
	Multiplication
	Division
)

func (o Operation) String() string {
	return operationsStr[o]
}
