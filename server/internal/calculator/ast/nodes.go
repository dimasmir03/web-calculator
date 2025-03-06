package ast

import (
	"github.com/dimasmir03/web-calculator-server/internal/calculator/lexer"
	"github.com/google/uuid"
	"github.com/m1gwings/treedrawer/tree"
)

func ToTreeDrawer(rootNode Node) *tree.Tree {
	t := tree.NewTree(nil)
	rootNode.toTreeDrawer(t)
	return t
}

type Node interface {
	toTreeDrawer(*tree.Tree)
	GetToken() *lexer.Token
}

var _ Node = &NumericNode{}
var _ Node = &UnaryNode{}
var _ Node = &BinaryNode{}

type NumericNode struct {
	val   float64
	token *lexer.Token
}

func NewNumericNode(val float64, token *lexer.Token) *NumericNode {
	return &NumericNode{
		val:   val,
		token: token,
	}
}

func (n *NumericNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeFloat64(n.val))
}
func (n *NumericNode) GetToken() *lexer.Token {
	return n.token
}
func (n *NumericNode) Value() float64 {
	return n.val
}

type UnaryNode struct {
	next     Node
	operator Operation
	token    *lexer.Token
}

func NewUnaryNode(operator Operation, next Node, token *lexer.Token) *UnaryNode {
	return &UnaryNode{
		operator: operator,
		next:     next,
		token:    token,
	}
}

func (n *UnaryNode) Operator() Operation {
	return n.operator
}
func (n *UnaryNode) Next() Node {
	return n.next
}
func (n *UnaryNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeString(n.operator.String()))
	n.next.toTreeDrawer(t.AddChild(nil))
}
func (n *UnaryNode) GetToken() *lexer.Token {
	return n.token
}

type UID string

type BinaryNode struct {
	uuid     UID
	operator Operation
	left     Node
	right    Node
	token    *lexer.Token
}

func NewBinaryNode(operator Operation, left, right Node, token *lexer.Token) *BinaryNode {
	return &BinaryNode{
		uuid:     UID(uuid.New().String()),
		operator: operator,
		left:     left,
		right:    right,
		token:    token,
	}
}

func (n *BinaryNode) Left() Node {
	return n.left
}
func (n *BinaryNode) Right() Node {
	return n.right
}
func (n *BinaryNode) Operator() Operation {
	return n.operator
}
func (n *BinaryNode) toTreeDrawer(t *tree.Tree) {
	t.SetVal(tree.NodeString(n.operator.String()))
	n.left.toTreeDrawer(t.AddChild(nil))
	n.right.toTreeDrawer(t.AddChild(nil))
}
func (n *BinaryNode) GetToken() *lexer.Token {
	return n.token
}
func (n *BinaryNode) GetUUID() UID {
	return n.uuid
}
