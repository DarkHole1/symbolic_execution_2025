package main

import (
	"symbolic-execution-course/internal/memory"
	"symbolic-execution-course/internal/symbolic"
)

func main() {
	var mem = memory.NewSymbolicMemory()
	var array = mem.Allocate(symbolic.ArrayType)

	mem.AssignToArray(array, 5, symbolic.NewIntConstant(10))

	var fromArray = mem.GetFromArray(array, 5)
	println(fromArray.String())

	// Will panic
	// var anotherFromArray = mem.GetFromArray(array, 10)
	// println(anotherFromArray)

	var object = mem.Allocate(symbolic.ObjectType)

	mem.AssignField(object, 5, symbolic.NewIntConstant(10))

	var fromObject = mem.GetFieldValue(object, 5)
	println(fromObject.String())
}
