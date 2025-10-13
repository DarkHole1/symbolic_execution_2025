// Package translator содержит реализацию транслятора в Z3
package translator

import (
	"fmt"
	"symbolic-execution-course/internal/symbolic"

	"github.com/ebukreev/go-z3/z3"
)

// Z3Translator транслирует символьные выражения в Z3 формулы
type Z3Translator struct {
	ctx    *z3.Context
	config *z3.Config
	vars   map[string]z3.Value // Кэш переменных
}

// NewZ3Translator создаёт новый экземпляр Z3 транслятора
func NewZ3Translator() *Z3Translator {
	config := &z3.Config{}
	ctx := z3.NewContext(config)

	return &Z3Translator{
		ctx:    ctx,
		config: config,
		vars:   make(map[string]z3.Value),
	}
}

// GetContext возвращает Z3 контекст
func (zt *Z3Translator) GetContext() interface{} {
	return zt.ctx
}

// Reset сбрасывает состояние транслятора
func (zt *Z3Translator) Reset() {
	zt.vars = make(map[string]z3.Value)
}

// Close освобождает ресурсы
func (zt *Z3Translator) Close() {
	// Z3 контекст закрывается автоматически
}

// TranslateExpression транслирует символьное выражение в Z3
func (zt *Z3Translator) TranslateExpression(expr symbolic.SymbolicExpression) (interface{}, error) {
	return expr.Accept(zt), nil
}

// VisitVariable транслирует символьную переменную в Z3
func (zt *Z3Translator) VisitVariable(expr *symbolic.SymbolicVariable) interface{} {
	v, ok := zt.vars[expr.Name]
	if ok {
		return v
	}

	v = zt.createZ3Variable(expr.Name, expr.ExprType)
	zt.vars[expr.Name] = v
	return v
}

// VisitIntConstant транслирует целочисленную константу в Z3
func (zt *Z3Translator) VisitIntConstant(expr *symbolic.IntConstant) interface{} {
	return zt.ctx.FromInt(expr.Value, zt.ctx.IntSort())
}

// VisitBoolConstant транслирует булеву константу в Z3
func (zt *Z3Translator) VisitBoolConstant(expr *symbolic.BoolConstant) interface{} {
	return zt.ctx.FromBool(expr.Value)
}

// VisitBinaryOperation транслирует бинарную операцию в Z3
func (zt *Z3Translator) VisitBinaryOperation(expr *symbolic.BinaryOperation) interface{} {
	left := expr.Left.Accept(zt).(z3.Int)
	right := expr.Right.Accept(zt).(z3.Int)

	switch expr.Operator {
	case symbolic.ADD:
		return left.Add(right)
	case symbolic.SUB:
		return left.Sub(right)
	case symbolic.DIV:
		return left.Div(right)
	case symbolic.MUL:
		return left.Mul(right)
	case symbolic.MOD:
		return left.Mod(right)
	case symbolic.EQ:
		return left.Eq(right)
	case symbolic.NE:
		return left.NE(right)
	case symbolic.LT:
		return left.LT(right)
	case symbolic.LE:
		return left.LE(right)
	case symbolic.GT:
		return left.GT(right)
	case symbolic.GE:
		return left.GE(right)
	}
	panic("not implemented")
}

// VisitLogicalOperation транслирует логическую операцию в Z3
func (zt *Z3Translator) VisitLogicalOperation(expr *symbolic.LogicalOperation) interface{} {
	operands := make([]z3.Bool, len(expr.Operands))
	for i, operand := range expr.Operands {
		operands[i] = operand.Accept(zt).(z3.Bool)
	}

	switch expr.Operator {
	case symbolic.AND:
		res := operands[0]
		for _, operand := range operands[1:] {
			res = res.And(operand)
		}
		return res
	case symbolic.OR:
		res := operands[0]
		for _, operand := range operands[1:] {
			res = res.Or(operand)
		}
		return res
	case symbolic.NOT:
		return operands[0].Not()
	case symbolic.IMPLIES:
		return operands[0].Implies(operands[1])
	}

	panic("not implemented")
}

// Вспомогательные методы

// createZ3Variable создаёт Z3 переменную соответствующего типа
func (zt *Z3Translator) createZ3Variable(name string, exprType symbolic.ExpressionType) z3.Value {
	switch exprType {
	case symbolic.IntType:
		return zt.ctx.IntConst(name)
	case symbolic.BoolType:
		return zt.ctx.BoolConst(name)
	case symbolic.ArrayType:
		return zt.ctx.ConstArray(zt.ctx.IntSort(), zt.ctx.IntConst(name))
	}
	panic("не реализовано")
}

// castToZ3Type приводит значение к нужному Z3 типу
func (zt *Z3Translator) castToZ3Type(value interface{}, targetType symbolic.ExpressionType) (z3.Value, error) {
	switch targetType {
	case symbolic.IntType:
		v, ok := value.(z3.Int)
		if !ok {
			return nil, fmt.Errorf("incorrect type cast")
		}
		return v, nil
	case symbolic.BoolType:
		v, ok := value.(z3.Bool)
		if !ok {
			return nil, fmt.Errorf("incorrect type cast")
		}
		return v, nil
	case symbolic.ArrayType:
		v, ok := value.(z3.Array)
		if !ok {
			return nil, fmt.Errorf("incorrect type cast")
		}
		return v, nil
	}

	panic("not implemented")
}
