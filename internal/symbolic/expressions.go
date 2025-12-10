// Package symbolic содержит конкретные реализации символьных выражений
package symbolic

import (
	"fmt"
	"strings"
)

// SymbolicExpression - базовый интерфейс для всех символьных выражений
type SymbolicExpression interface {
	// Type возвращает тип выражения
	Type() ExpressionType

	// String возвращает строковое представление выражения
	String() string

	// Accept принимает visitor для обхода дерева выражений
	Accept(visitor Visitor) interface{}
}

// SymbolicVariable представляет символьную переменную
type SymbolicVariable struct {
	Name     string
	ExprType ExpressionType
}

// NewSymbolicVariable создаёт новую символьную переменную
func NewSymbolicVariable(name string, exprType ExpressionType) *SymbolicVariable {
	return &SymbolicVariable{
		Name:     name,
		ExprType: exprType,
	}
}

// Type возвращает тип переменной
func (sv *SymbolicVariable) Type() ExpressionType {
	return sv.ExprType
}

// String возвращает строковое представление переменной
func (sv *SymbolicVariable) String() string {
	return sv.Name
}

// Accept реализует Visitor pattern
func (sv *SymbolicVariable) Accept(visitor Visitor) interface{} {
	return visitor.VisitVariable(sv)
}

// IntConstant представляет целочисленную константу
type IntConstant struct {
	Value int64
}

// NewIntConstant создаёт новую целочисленную константу
func NewIntConstant(value int64) *IntConstant {
	return &IntConstant{Value: value}
}

// Type возвращает тип константы
func (ic *IntConstant) Type() ExpressionType {
	return IntType
}

// String возвращает строковое представление константы
func (ic *IntConstant) String() string {
	return fmt.Sprintf("%d", ic.Value)
}

// Accept реализует Visitor pattern
func (ic *IntConstant) Accept(visitor Visitor) interface{} {
	return visitor.VisitIntConstant(ic)
}

// BoolConstant представляет булеву константу
type BoolConstant struct {
	Value bool
}

// NewBoolConstant создаёт новую булеву константу
func NewBoolConstant(value bool) *BoolConstant {
	return &BoolConstant{Value: value}
}

// Type возвращает тип константы
func (bc *BoolConstant) Type() ExpressionType {
	return BoolType
}

// String возвращает строковое представление константы
func (bc *BoolConstant) String() string {
	return fmt.Sprintf("%t", bc.Value)
}

// Accept реализует Visitor pattern
func (bc *BoolConstant) Accept(visitor Visitor) interface{} {
	return visitor.VisitBoolConstant(bc)
}

type FloatConstant struct {
	Value float64
}

// NewBoolConstant создаёт новую булеву константу
func NewFloatConstant(value float64) *FloatConstant {
	return &FloatConstant{Value: value}
}

// Type возвращает тип константы
func (fc *FloatConstant) Type() ExpressionType {
	return FloatType
}

// String возвращает строковое представление константы
func (fc *FloatConstant) String() string {
	return fmt.Sprintf("%f", fc.Value)
}

// Accept реализует Visitor pattern
func (fc *FloatConstant) Accept(visitor Visitor) interface{} {
	return visitor.VisitFloatConstant(fc)
}

// BinaryOperation представляет бинарную операцию
type BinaryOperation struct {
	Left     SymbolicExpression
	Right    SymbolicExpression
	Operator BinaryOperator
}

// NewBinaryOperation создаёт новую бинарную операцию
func NewBinaryOperation(left, right SymbolicExpression, op BinaryOperator) *BinaryOperation {
	if left.Type() == BoolType && right.Type() == BoolType {
		return &BinaryOperation{
			Left:     left,
			Right:    right,
			Operator: op,
		}
	}

	if left.Type() == IntType && right.Type() == IntType {
		return &BinaryOperation{
			Left:     left,
			Right:    right,
			Operator: op,
		}
	}

	if left.Type() == FloatType && right.Type() == FloatType {
		return &BinaryOperation{
			Left:     left,
			Right:    right,
			Operator: op,
		}
	}

	if left.Type() == IntType && right.Type() == FloatType || left.Type() == FloatType && right.Type() == IntType {
		return &BinaryOperation{
			Left:     left,
			Right:    right,
			Operator: op,
		}
	}

	panic("incompatible types")
}

// Type возвращает результирующий тип операции
func (bo *BinaryOperation) Type() ExpressionType {
	switch bo.Operator {
	case ADD, SUB, MUL, DIV, MOD, BAND, BOR, BXOR, SHL, SHR:
		return IntType
	case EQ, NE, GT, LT, GE, LE:
		return BoolType
	}
	panic("not implemented")
}

// String возвращает строковое представление операции
func (bo *BinaryOperation) String() string {
	return fmt.Sprintf("(%s %s %s)", bo.Left.String(), bo.Operator.String(), bo.Right.String())
}

// Accept реализует Visitor pattern
func (bo *BinaryOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryOperation(bo)
}

// LogicalOperation представляет логическую операцию
type LogicalOperation struct {
	Operands []SymbolicExpression
	Operator LogicalOperator
}

// NewLogicalOperation создаёт новую логическую операцию
func NewLogicalOperation(operands []SymbolicExpression, op LogicalOperator) *LogicalOperation {
	switch op {
	case AND, OR:
		if len(operands) < 2 {
			panic("incorrect number of arguments")
		}
	case IMPLIES:
		if len(operands) != 2 {
			panic("incorrect number of arguments")
		}
	case NOT:
		if len(operands) != 1 {
			panic("incorrect number of arguments")
		}
	}
	for _, op := range operands {
		if op.Type() != BoolType {
			panic("incorrect type")
		}
	}
	return &LogicalOperation{
		Operands: operands,
		Operator: op,
	}
}

// Type возвращает тип логической операции (всегда bool)
func (lo *LogicalOperation) Type() ExpressionType {
	return BoolType
}

// String возвращает строковое представление логической операции
func (lo *LogicalOperation) String() string {
	// TODO: Реализовать
	// Для NOT: "!operand"
	// Для AND/OR: "(operand1 && operand2 && ...)"
	// Для IMPLIES: "(operand1 => operand2)"
	switch lo.Operator {
	case NOT:
		return fmt.Sprintf("!%s", lo.Operands[0].String())
	case AND, OR:
		strOperands := make([]string, len(lo.Operands))
		for i, operand := range lo.Operands {
			strOperands[i] = operand.String()
		}
		return fmt.Sprintf("(%s)", strings.Join(strOperands, " "+lo.Operator.String()+" "))
	case IMPLIES:
		return fmt.Sprintf("(%s => %s)", lo.Operands[0].String(), lo.Operands[1].String())
	}
	panic("not implemented")
}

// Accept реализует Visitor pattern
func (lo *LogicalOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitLogicalOperation(lo)
}

// Операторы для бинарных выражений
type BinaryOperator int

const (
	// Арифметические операторы
	ADD BinaryOperator = iota
	SUB
	MUL
	DIV
	MOD

	// Операторы сравнения
	EQ // равно
	NE // не равно
	LT // меньше
	LE // меньше или равно
	GT // больше
	GE // больше или равно

	BAND
	BOR
	BXOR
	SHL
	SHR
)

// String возвращает строковое представление оператора
func (op BinaryOperator) String() string {
	switch op {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	case MOD:
		return "%"
	case EQ:
		return "=="
	case NE:
		return "!="
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	default:
		return "unknown"
	}
}

// Логические операторы
type LogicalOperator int

const (
	AND LogicalOperator = iota
	OR
	NOT
	IMPLIES
)

// String возвращает строковое представление логического оператора
func (op LogicalOperator) String() string {
	switch op {
	case AND:
		return "&&"
	case OR:
		return "||"
	case NOT:
		return "!"
	case IMPLIES:
		return "=>"
	default:
		return "unknown"
	}
}

type UnaryOperator int

const (
	BNOT UnaryOperator = iota
)

func (op UnaryOperator) String() string {
	switch op {
	case BNOT:
		return "^"
	default:
		return "unknown"
	}
}

type Ref struct {
	Tpe    ExpressionType
	Ptr    int64
	Memory interface {
		Deref(*Ref) SymbolicExpression
	}
}

func (ref *Ref) Type() ExpressionType {
	return ReferenceType
}

func (ref *Ref) String() string {
	return fmt.Sprintf("0x%04x", ref.Ptr)
}

func (ref *Ref) Accept(visitor Visitor) interface{} {
	return ref.Memory.Deref(ref).Accept(visitor)
}

type UnaryOperation struct {
	Left     SymbolicExpression
	Operator UnaryOperator
}

func NewUnaryOperation(left SymbolicExpression, op UnaryOperator) *UnaryOperation {
	if left.Type() != IntType {
		panic("incompatible types")
	}

	return &UnaryOperation{
		Left:     left,
		Operator: op,
	}
}

// Type возвращает результирующий тип операции
func (uo *UnaryOperation) Type() ExpressionType {
	return IntType
}

// String возвращает строковое представление операции
func (uo *UnaryOperation) String() string {
	return fmt.Sprintf("(%s %s)", uo.Operator.String(), uo.Left.String())
}

// Accept реализует Visitor pattern
func (uo *UnaryOperation) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryOperation(uo)
}

// TODO: Добавьте дополнительные типы выражений по необходимости:
// - ArrayAccess (доступ к элементам массива: arr[index])
// - FunctionCall (вызовы функций: f(x, y))
// - ConditionalExpression (тернарный оператор: condition ? true_expr : false_expr)
