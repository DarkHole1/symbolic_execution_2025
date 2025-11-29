package internal

import (
	"fmt"
	"go/constant"
	"go/types"
	"symbolic-execution-course/internal/memory"
	"symbolic-execution-course/internal/symbolic"

	"golang.org/x/tools/go/ssa"
)

type Interpreter struct {
	CallStack     []CallStackFrame
	Analyser      *Analyser
	PathCondition symbolic.SymbolicExpression
	Heap          memory.Memory
}

type CallStackFrame struct {
	Function     *ssa.Function
	LocalMemory  map[string]symbolic.SymbolicExpression
	ReturnValue  symbolic.SymbolicExpression
	CurrentBlock int
	CurrentInstr int
}

func ConvertType(tpe types.Type) symbolic.ExpressionType {
	switch tpe.Underlying().(*types.Basic).Kind() {
	case types.Bool:
		return symbolic.BoolType
	case types.Int:
		return symbolic.IntType
	default:
		panic("unexpected types.BasicKind")
	}
}

func (interpreter *Interpreter) interpretDynamically(element ssa.Instruction) []Interpreter {
	switch element := element.(type) {
	case *ssa.BinOp:
		X := interpreter.resolveExpression(element.X)
		Y := interpreter.resolveExpression(element.Y)
		fmt.Println("binop", element.Op, X, Y)

	case *ssa.Alloc:
	case *ssa.Call:
	case *ssa.ChangeInterface:
	case *ssa.ChangeType:
	case *ssa.Convert:
	case *ssa.DebugRef:
	case *ssa.Defer:
	case *ssa.Extract:
	case *ssa.Field:
	case *ssa.FieldAddr:
	case *ssa.Go:
	case *ssa.If:
	case *ssa.Index:
	case *ssa.IndexAddr:
	case *ssa.Jump:
	case *ssa.Lookup:
	case *ssa.MakeChan:
	case *ssa.MakeClosure:
	case *ssa.MakeInterface:
	case *ssa.MakeMap:
	case *ssa.MakeSlice:
	case *ssa.MapUpdate:
	case *ssa.MultiConvert:
	case *ssa.Next:
	case *ssa.Panic:
	case *ssa.Phi:
	case *ssa.Range:
	case *ssa.Return:
	case *ssa.RunDefers:
	case *ssa.Select:
	case *ssa.Send:
	case *ssa.Slice:
	case *ssa.SliceToArrayPointer:
	case *ssa.Store:
	case *ssa.TypeAssert:
	case *ssa.UnOp:
	default:
		panic(fmt.Sprintf("unexpected ssa.Instruction: %#v", element))
	}
	return []Interpreter{*interpreter}
}

func (interpreter *Interpreter) resolveExpression(value ssa.Value) symbolic.SymbolicExpression {
	switch value := value.(type) {
	case *ssa.Const:
		switch value.Type().Underlying().(*types.Basic).Kind() {
		case types.Int:
			return symbolic.NewIntConstant(value.Int64())
		case types.Bool:
			return symbolic.NewBoolConstant(constant.BoolVal(value.Value))
		default:
			panic(fmt.Sprintf("unexpected value.Kind(): %#v", value.Type().Underlying().(*types.Basic).Kind()))
		}

	case *ssa.Parameter:
		frame := interpreter.CallStack[len(interpreter.CallStack)-1]
		res := frame.LocalMemory[value.Name()]
		if res == nil {
			panic(fmt.Sprintf("no parameter in scope %#v", value.Name()))
		}
		return res

	default:
		panic(fmt.Sprintf("unexpected ssa.Value: %#v", value))
	}
}
