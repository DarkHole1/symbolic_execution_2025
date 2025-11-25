package memory

import "symbolic-execution-course/internal/symbolic"

type Memory interface {
	Allocate(tpe symbolic.ExpressionType) *symbolic.Ref

	AssignField(ref *symbolic.Ref, fieldIdx int, value symbolic.SymbolicExpression)

	GetFieldValue(ref *symbolic.Ref, fieldIdx int) symbolic.SymbolicExpression

	AssignToArray(ref *symbolic.Ref, index int, value symbolic.SymbolicExpression)

	GetFromArray(ref *symbolic.Ref, index int) symbolic.SymbolicExpression
}

type SymbolicMemory struct {
	c               int64
	pool            map[symbolic.ExpressionType]map[int64]symbolic.SymbolicExpression
	arrayObjectPool map[int64]map[int]symbolic.SymbolicExpression
}

func NewSymbolicMemory() SymbolicMemory {
	return SymbolicMemory{
		c:               0,
		pool:            make(map[symbolic.ExpressionType]map[int64]symbolic.SymbolicExpression),
		arrayObjectPool: make(map[int64]map[int]symbolic.SymbolicExpression),
	}
}

func (mem *SymbolicMemory) Allocate(tpe symbolic.ExpressionType) *symbolic.Ref {
	mem.c += 1
	switch tpe {
	case symbolic.ArrayType:
		mem.arrayObjectPool[mem.c] = make(map[int]symbolic.SymbolicExpression)
	case symbolic.ObjectType:
		mem.arrayObjectPool[mem.c] = make(map[int]symbolic.SymbolicExpression)
	default:
		if mem.pool[tpe] == nil {
			mem.pool[tpe] = make(map[int64]symbolic.SymbolicExpression)
		}
	}
	return &symbolic.Ref{
		Tpe: tpe,
		Ptr: mem.c,
	}
}

func (mem *SymbolicMemory) AssignField(ref *symbolic.Ref, fieldIdx int, value symbolic.SymbolicExpression) {
	if ref.Tpe != symbolic.ObjectType {
		panic("incorrect type")
	}

	mem.arrayObjectPool[ref.Ptr][fieldIdx] = value
}

func (mem *SymbolicMemory) GetFieldValue(ref *symbolic.Ref, fieldIdx int) symbolic.SymbolicExpression {
	if ref.Tpe != symbolic.ObjectType {
		panic("incorrect type")
	}

	value, ok := mem.arrayObjectPool[ref.Ptr][fieldIdx]
	if !ok {
		panic("undefined object field")
	}

	return value
}

func (mem *SymbolicMemory) AssignToArray(ref *symbolic.Ref, index int, value symbolic.SymbolicExpression) {
	if ref.Tpe != symbolic.ArrayType {
		panic("incorrect type")
	}

	mem.arrayObjectPool[ref.Ptr][index] = value
}

func (mem *SymbolicMemory) GetFromArray(ref *symbolic.Ref, index int) symbolic.SymbolicExpression {
	if ref.Tpe != symbolic.ArrayType {
		panic("incorrect type")
	}

	value, ok := mem.arrayObjectPool[ref.Ptr][index]
	if !ok {
		panic("undefined array index")
	}

	return value
}
