package calculator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/dimasmir03/web-calculator-server/internal/calculator/ast"
	"github.com/dimasmir03/web-calculator-server/internal/calculator/evaluator"
	lexer3 "github.com/dimasmir03/web-calculator-server/internal/calculator/lexer"
	parser3 "github.com/dimasmir03/web-calculator-server/internal/calculator/parser"
	"github.com/dimasmir03/web-calculator-server/internal/calculator/parser/recursivedescent"
	"github.com/dimasmir03/web-calculator-server/internal/model"
	"github.com/dimasmir03/web-calculator-server/internal/storage/sqlite"
)

// func main() {
// 	calc := NewCalculator()
// 	var s string
// 	// fmt.Scan(&s)
// 	// s = "1*2+3*(4+5)+5"
// 	s = "1+2*3"
// 	fmt.Println(calc.AddExpr(s))
// 	s = "1*2+3*3*(3+4)+3*34"
// 	fmt.Println(calc.AddExpr(s))
// 	calc.ShowResults()
// 	// fmt.Println(calc.queue)
// }

func NewCalculator(db *sqlite.Storage) *Calculator {
	return &Calculator{
		expr: make([]*Expression, 0),
		db:   db,
	}
}

type SimpleExpression struct {
	A      any
	B      any
	Op     string
	Id     ast.UID
	status string
}

type Expression struct {
	id         ast.UID
	expression string
	parser     parser3.Parser
	evaluator  *evaluator.NumericEvaluator
	lexer      *lexer3.Lexer
	tokenized  []*lexer3.Token
	rootNode   ast.Node
	value      float64
	status     string
	queue      map[ast.UID]*SimpleExpression
}

type Calculator struct {
	expr []*Expression
	db   *sqlite.Storage
	// queue []*SimpleExpression
	m sync.Mutex
}

type Expr struct {
	Id         ast.UID `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
}

func (c *Calculator) GetExpressionById(reqid string) *Expr {
	res := new(Expr)
	id := ast.UID(reqid)
	for _, v := range c.expr {
		if v.id == id {
			res.Expression = v.expression
			res.Id = v.id
			res.Result = v.value
			res.Status = v.status
			return res
		}
	}
	return nil
}

func (c *Calculator) GetExpressionsStatus() []Expr {
	var result []Expr
	for _, v := range c.expr {
		result = append(result, Expr{
			Id:         v.id,
			Expression: v.expression,
			Status:     v.status,
			Result:     v.value,
		})
	}
	return result
}

func (c *Calculator) AddExpr(expr string) (string, error) {
	exp := &Expression{
		expression: expr,
		queue:      make(map[ast.UID]*SimpleExpression),
		status:     "pending",
	}
	c.m.Lock()
	defer c.m.Unlock()
	id, err := exp.Eval()
	if err != nil {
		return "", err
	}
	c.expr = append(c.expr, exp)
	// c.queue = append(c.queue, ...)
	// fmt.Println("queue", c.queue)
	return id, nil
}

func (c *Expression) Eval() (string, error) {
	c.expression = strings.TrimSpace(c.expression)
	var err error
	c.parser, err = recursivedescent.NewParser(parser3.DefaultTokenPriorities())
	if err != nil {
		return "", err
	}
	c.evaluator, err = evaluator.NewNumericEvaluator()
	if err != nil {
		return "", err
	}
	c.lexer = lexer3.NewLexer(c.expression)
	c.tokenized, err = c.lexer.Tokenize()
	if err != nil {
		return "", err
	}
	c.rootNode, err = c.parser.Parse(c.tokenized)
	if err != nil {
		return "", err
	}
	c.id = c.rootNode.GetUUID()
	if n, ok := c.rootNode.(*ast.UnaryNode); ok {
		res, _ := c.evaluator.HandleUnary(n)
		c.id = n.GetUUID()
		c.value = res
		c.status = "complete"
	} else {
		c.queueFromNode(c.rootNode)
	}
	// c.sortQueue()
	// c.value, err = c.evaluator.Eval(c.rootNode)
	// if err != nil {
	// 	return "", err
	// }
	return string(c.id), nil
}

func (c *Expression) queueFromNode(rootNode ast.Node) any {
	switch n := rootNode.(type) {
	case *ast.BinaryNode:
		c.queue[n.GetUUID()] =
			&SimpleExpression{c.queueFromNode(n.Left()),
				c.queueFromNode(n.Right()),
				n.Operator().String(),
				n.GetUUID(),
				"open"}
		// c.queue = append(c.queue,
		// 	&SimpleExpression{fmt.Sprint(c.queueFromNode(n.Left())),
		// 		fmt.Sprint(c.queueFromNode(n.Right())),
		// 		n.GetToken().Type().String(),
		// 		n.GetUUID(),
		// 		"open"})
		return n.GetUUID()
	case *ast.UnaryNode:
		res, _ := c.evaluator.HandleUnary(n)
		return res
	case *ast.NumericNode:
		return n.Value()
	}
	return ""
}

func (c *Calculator) ShowResults() {
	printTree := true
	c.m.Lock()
	defer c.m.Unlock()
	for _, v := range c.expr {
		// fmt.Println(v.queue)
		for _, v := range v.queue {
			fmt.Print(v)
		}
		fmt.Println()
		if printTree {
			fmt.Print(ast.ToTreeDrawer(v.rootNode))
		}
	}
}

func (c *Calculator) GetSimpleExpr() SimpleExpression {
	var res SimpleExpression
	for _, v := range c.expr {
		fmt.Println(&v, v)
	}
	for i := range c.expr {
		if strings.HasPrefix(c.expr[i].status, "error") {
			continue
		}
		for uid, v := range c.expr[i].queue {
			fmt.Println(v)
			fmt.Println("and here", reflect.TypeOf(v.A), reflect.TypeOf(v.B), v.status == "open")
			if isNumber(v.A) &&
				isNumber(v.B) &&
				v.status == "open" {
				c.expr[i].queue[uid].status = "block"
				go func(expression *SimpleExpression) {
					timer := time.NewTimer(10 * time.Second)
					<-timer.C
					expression.status = "open"
				}(c.expr[i].queue[uid])
				fmt.Println("there")
				res.A = c.expr[i].queue[uid].A
				res.Id = c.expr[i].queue[uid].Id
				res.B = c.expr[i].queue[uid].B
				res.status = c.expr[i].queue[uid].status
				res.Op = c.expr[i].queue[uid].Op
				return res
			}
		}
	}
	fmt.Println("here")
	return SimpleExpression{}
}

func (c *Calculator) SetSimpleExprResult(nid string, result float64, error string) error {
	c.m.Lock()
	defer c.m.Unlock()
	id := ast.UID(nid)
	var flag bool
	for i := range c.expr {
		if strings.HasPrefix(c.expr[i].status, "error") {
			continue
		}
		for j, v := range c.expr[i].queue {
			if v.Id == id {
				if error != "" {
					c.expr[i].status = "error: " + error
					return nil
					// c.expr[i]
				}
				if len(c.expr[i].queue) == 1 {
					c.expr[i].status = "complete"
					c.db.UpdateExpression(&model.Expression{ID: string(c.expr[i].id), Result: result, Status: "complete"})
					c.expr[i].value = result
				}
				delete(c.expr[i].queue, id)
				flag = true
			} else {
				switch v.A.(type) {
				case ast.UID:
					if v.A.(ast.UID) == id {
						c.expr[i].queue[j].A = result
						flag = true
					}
				}
				switch v.B.(type) {
				case ast.UID:
					if v.B.(ast.UID) == id {
						c.expr[i].queue[j].B = result
						flag = true
					}
				}
			}
		}
	}
	if !flag {
		return errors.New("not found by")
	}
	return nil
}

// func (c *Expression) sortQueue() []*SimpleExpression {
// 	var numbersQueue, uuidQueue []*SimpleExpression
// 	for _, item := range c.queue {
// 		if isNumber(item.A) && isNumber(item.B) {
// 			numbersQueue = append(numbersQueue, item)
// 			fmt.Println("numbers1", numbersQueue)
// 		} else {
// 			uuidQueue = append(uuidQueue, item)
// 		}
// 	}
// 	fmt.Println("numbers1", numbersQueue)
//
// 	c.queue = append(numbersQueue, uuidQueue...)
// 	fmt.Println("numbers2", c.queue)
// 	fmt.Println("numbers2", numbersQueue)
// 	return numbersQueue
// }

func isNumber(x interface{}) bool {
	_, ok := x.(float64)
	return ok
}
