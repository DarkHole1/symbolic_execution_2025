package internal

import (
	"symbolic-execution-course/internal/memory"
	issa "symbolic-execution-course/internal/ssa"
	"symbolic-execution-course/internal/symbolic"
	"symbolic-execution-course/internal/translator"

	"golang.org/x/tools/go/ssa"
)

type Analyser struct {
	Package      *ssa.Package
	StatesQueue  PriorityQueue
	PathSelector PathSelector
	Results      []Interpreter
	Z3Translator *translator.Z3Translator
}

func createAnalyser(source string, functionName string, selector PathSelector) *Analyser {
	builder := issa.NewBuilder()

	graph, err := builder.ParseAndBuildSSA(source, functionName)
	if err != nil {
		panic("ssa parsing failed")
	}

	zt := translator.NewZ3Translator()

	frame := CallStackFrame{
		Function:     graph,
		LocalMemory:  map[ssa.Value]symbolic.SymbolicExpression{},
		ReturnValue:  nil,
		CurrentBlock: 0,
	}

	for _, param := range graph.Params {
		frame.LocalMemory[param] = symbolic.NewSymbolicVariable(param.Name(), ConvertType(param.Type()))
	}

	res := &Analyser{
		Package:      graph.Package(),
		PathSelector: selector,
		Results:      []Interpreter{},
		Z3Translator: zt,
	}

	start := Interpreter{
		CallStack: []CallStackFrame{
			frame,
		},
		PathCondition: symbolic.NewBoolConstant(true),
		Heap:          memory.NewSymbolicMemory(),
		Analyser:      res,
	}

	var queue PriorityQueue

	queue.Push(&Item{
		value:    start,
		priority: selector.CalculatePriority(start),
	})

	res.StatesQueue = queue

	return res
}

func Analyse(source string, functionName string) []Interpreter {
	analyser := createAnalyser(source, functionName, &RandomPathSelector{})

	i := 0
	for i < 10 && analyser.StatesQueue.Len() > 0 {
		state := analyser.StatesQueue.Pop().(*Item).value
		new_states := state.interpretCurrentBlock()
		for _, new_state := range new_states {
			analyser.StatesQueue.Push(&Item{
				value:    new_state,
				priority: analyser.PathSelector.CalculatePriority(new_state),
			})
		}
		i++
	}

	return analyser.Results
}
