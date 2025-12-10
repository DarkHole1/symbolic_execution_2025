package main

import (
	"fmt"
	"os"
	"symbolic-execution-course/internal"
)

func main() {
	source_bytes, _ := os.ReadFile("./examples/test_functions.go")

	source := string(source_bytes)
	test_functions := []string{
		"test1",
		"test2",
		"testArithmetic",
		"testUnary",
		"testComparisons",
		"testLogicalOps",
		"testWhileLoop",
		"testForLoop",
		"testInfiniteLoopBreak",
		"testLoopWithConcreteBoundAndSymbolicBranching",
		"testLoopWithSymbolicBoundAndSymbolicBranching",
		"testComplexConditions",
		"testNestedIf",
		"testBitwise",
		"testCombined",
		"testEdgeCases",
		"testMultipleReturns",
		// "testSimpleSum",
	}

	for _, fun := range test_functions {
		fmt.Printf("=== %s ===\n", fun)
		result := internal.Analyse(string(source), fun)
		for _, interpreter := range result {
			fmt.Println(interpreter)
		}
	}

}
