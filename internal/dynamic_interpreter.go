package internal

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"slices"
	"symbolic-execution-course/internal/memory"
	"symbolic-execution-course/internal/symbolic"

	"github.com/LastPossum/kamino"
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
	LocalMemory  map[ssa.Value]symbolic.SymbolicExpression
	ReturnValue  []symbolic.SymbolicExpression
	CurrentBlock int
	PrevBlock    int
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

func execBinOp(op token.Token, X, Y symbolic.SymbolicExpression) symbolic.SymbolicExpression {
	switch op {
	case token.ADD:
		return symbolic.NewBinaryOperation(X, Y, symbolic.ADD)
	case token.SUB:
		return symbolic.NewBinaryOperation(X, Y, symbolic.SUB)
	case token.MUL:
		return symbolic.NewBinaryOperation(X, Y, symbolic.MUL)
	case token.QUO:
		return symbolic.NewBinaryOperation(X, Y, symbolic.DIV)
	case token.REM:
		return symbolic.NewBinaryOperation(X, Y, symbolic.MOD)
	case token.EQL:
		return symbolic.NewBinaryOperation(X, Y, symbolic.EQ)
	case token.LSS:
		return symbolic.NewBinaryOperation(X, Y, symbolic.LT)
	case token.GTR:
		return symbolic.NewBinaryOperation(X, Y, symbolic.GT)
	case token.NEQ:
		return symbolic.NewBinaryOperation(X, Y, symbolic.NE)
	case token.LEQ:
		return symbolic.NewBinaryOperation(X, Y, symbolic.LE)
	case token.GEQ:
		return symbolic.NewBinaryOperation(X, Y, symbolic.GE)
	case token.LAND:
		return symbolic.NewLogicalOperation([]symbolic.SymbolicExpression{X, Y}, symbolic.AND)
	case token.LOR:
		return symbolic.NewLogicalOperation([]symbolic.SymbolicExpression{X, Y}, symbolic.OR)
	default:
		panic(fmt.Sprintf("unexpected token.Token: %#v", op))
	}
}

func execUnOp(op token.Token, X symbolic.SymbolicExpression) symbolic.SymbolicExpression {
	switch op {
	case token.SUB:
		return symbolic.NewBinaryOperation(X, symbolic.NewIntConstant(-1), symbolic.MUL)
	case token.NOT:
		return symbolic.NewLogicalOperation([]symbolic.SymbolicExpression{X}, symbolic.NOT)
	default:
		panic(fmt.Sprintf("unexpected token.Token: %#v", op))
	}
}

func (interpreter *Interpreter) frame() *CallStackFrame {
	return &interpreter.CallStack[len(interpreter.CallStack)-1]
}

func (interpreter *Interpreter) interpretCurrentBlock() []Interpreter {
	var res []Interpreter
	nonPhis := interpreter.executePhis()
	for _, instr := range nonPhis {
		res = interpreter.interpretDynamically(instr)
	}
	return res
}

func (interpreter *Interpreter) executePhis() []ssa.Instruction {
	frame := interpreter.frame()
	block := frame.Function.Blocks[frame.CurrentBlock]

	firstNonPhi := -1
	for i, instr := range block.Instrs {
		if _, ok := instr.(*ssa.Phi); !ok {
			firstNonPhi = i
			break
		}
	}

	nonPhis := block.Instrs[firstNonPhi:]
	if firstNonPhi > 0 {
		phis := block.Instrs[:firstNonPhi]

		predIndex := slices.Index(block.Preds, frame.Function.Blocks[frame.PrevBlock])
		phitemps := make([]symbolic.SymbolicExpression, len(phis))
		for i, phi := range phis {
			phi := phi.(*ssa.Phi)
			edge := phi.Edges[predIndex]
			phitemps[i] = interpreter.resolveExpression(edge)
		}
		for i, phi := range phis {
			interpreter.frame().LocalMemory[phi.(*ssa.Phi)] = phitemps[i]
		}
	}
	return nonPhis
}

func (interpreter *Interpreter) interpretDynamically(element ssa.Instruction) []Interpreter {
	switch element := element.(type) {
	case *ssa.BinOp:
		X := interpreter.resolveExpression(element.X)
		Y := interpreter.resolveExpression(element.Y)
		interpreter.frame().LocalMemory[element] = execBinOp(element.Op, X, Y)
		return nil

	case *ssa.UnOp:
		X := interpreter.resolveExpression(element.X)
		interpreter.frame().LocalMemory[element] = execUnOp(element.Op, X)
		return nil

	case *ssa.If:
		cond := interpreter.resolveExpression(element.Cond)
		intTrue, _ := kamino.Clone(interpreter)
		intFalse, _ := kamino.Clone(interpreter)

		// TODO: Fix logic
		intTrue.Analyser = interpreter.Analyser
		intFalse.Analyser = interpreter.Analyser

		succs := interpreter.frame().Function.Blocks[interpreter.frame().CurrentBlock].Succs

		intTrue.PathCondition = symbolic.NewLogicalOperation(
			[]symbolic.SymbolicExpression{intTrue.PathCondition, cond},
			symbolic.AND,
		)
		intTrue.frame().CurrentBlock = succs[0].Index
		intTrue.frame().PrevBlock = interpreter.frame().CurrentBlock
		intFalse.PathCondition = symbolic.NewLogicalOperation(
			[]symbolic.SymbolicExpression{
				intFalse.PathCondition,
				symbolic.NewLogicalOperation([]symbolic.SymbolicExpression{cond}, symbolic.NOT),
			},
			symbolic.AND,
		)
		intFalse.frame().CurrentBlock = succs[1].Index
		intFalse.frame().PrevBlock = interpreter.frame().CurrentBlock
		return []Interpreter{*intTrue, *intFalse}

		// case *ssa.Alloc:
	case *ssa.Jump:
		interpreter.frame().PrevBlock = interpreter.frame().CurrentBlock
		interpreter.frame().CurrentBlock = element.Block().Succs[0].Index
		return []Interpreter{*interpreter}

	case *ssa.Return:
		results := make([]symbolic.SymbolicExpression, len(element.Results))
		for i, res := range element.Results {
			results[i] = interpreter.resolveExpression(res)
		}
		interpreter.frame().ReturnValue = results
		interpreter.Analyser.Results = append(interpreter.Analyser.Results, *interpreter)
		return []Interpreter{}
	default:
		panic(fmt.Sprintf("unexpected ssa.Instruction: %#v", element))
	}
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

	default:
		if v, ok := interpreter.frame().LocalMemory[value]; ok {
			return v
		}
		panic(fmt.Sprintf("unexpected ssa.Value: %#v", value))
	}
}
