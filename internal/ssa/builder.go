// Package ssa предоставляет функции для построения SSA представления
package ssa

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Builder отвечает за построение SSA из исходного кода Go
type Builder struct {
	fset *token.FileSet
}

// NewBuilder создаёт новый экземпляр Builder
func NewBuilder() *Builder {
	return &Builder{
		fset: token.NewFileSet(),
	}
}

// ParseAndBuildSSA парсит исходный код Go и создаёт SSA представление
// Возвращает SSA программу и функцию по имени
func (b *Builder) ParseAndBuildSSA(source string, funcName string) (*ssa.Function, error) {
	file, err := parser.ParseFile(b.fset, "test.go", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	files := []*ast.File{file}

	pkg := types.NewPackage(file.Name.Name, "")

	res, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()}, b.fset, pkg, files, ssa.SanityCheckFunctions)
	if err != nil {
		return nil, err
	}

	f := res.Func(funcName)

	if f == nil {
		return nil, fmt.Errorf("no such function")
	}

	return f, nil
}
