package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

// processFile 处理单个Go文件
func processFile(path string) {
	fset := token.NewFileSet()

	// 解析文件
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", path, err)
		return
	}

	identName := IDENT_NAME
	for _, im := range f.Imports {
		if im.Path.Value == "\""+packageName+"\"" {
			if im.Name != nil {
				identName = im.Name.Name
			}
			break
		}
	}

	// 查找并替换errorn.Wrap调用
	modified := false
	// 函数声明栈（用于追踪当前函数上下文）
	var funcStack []ast.Node

	preFunc := func(cursor *astutil.Cursor) bool {
		n := cursor.Node()

		// 压入函数声明节点
		switch n.(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			funcStack = append(funcStack, n)
		}

		es, ok := n.(*ast.ExprStmt)
		if !ok {
			return true
		}

		call, ok := es.X.(*ast.CallExpr)
		if !ok {
			return true
		}

		// 检查是否是errorn.Wrap调用
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if ident, ok := sel.X.(*ast.Ident); !ok || ident.Name != identName {
			return true
		}

		flag, ok := FLAG_FUNC_NAME_MAP[sel.Sel.Name]
		if !ok {
			return true
		}

		// 获取位置信息
		pos := fset.Position(call.Pos())

		// 获取当前函数上下文
		var currentFunc ast.Node
		if len(funcStack) > 0 {
			currentFunc = funcStack[len(funcStack)-1]
		}

		// 提取参数
		if len(call.Args) == 0 {
			fmt.Printf("Warning: errorn.Wrap called without arguments at %s: %s\n", path, pos)
			return true
		}

		errIdent, ok := call.Args[0].(*ast.Ident)
		if !ok {
			fmt.Printf("Warning: first argument to errorn.Wrap must be identifier at %s: %s\n", path, pos)
			return true
		}

		// 根据函数返回类型生成替换语句
		var replacement ast.Stmt
		if currentFunc != nil {
			returnInfo := getReturnInfo(currentFunc)
			if err := validateReturnTypes(returnInfo.types); err != nil {
				fmt.Printf("Warning: %s at %s\n", err, pos)
				return true
			}
			replacement = generateReplacement(flag, identName, errIdent, call.Args[1:], returnInfo)
		} else {
			// 没有函数上下文时使用默认处理
			replacement = generateDefault(flag, identName, errIdent, call.Args[1:])
		}

		// 替换
		cursor.Replace(replacement)
		modified = true
		return false
	}

	postFunc := func(cursor *astutil.Cursor) bool {
		switch cursor.Node().(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			if len(funcStack) > 0 {
				funcStack = funcStack[:len(funcStack)-1]
			}
		}
		return true
	}

	astutil.Apply(f, preFunc, postFunc)

	// 如果文件有修改，保存
	if modified {
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, f); err != nil {
			fmt.Printf("Error formatting %s: %v\n", path, err)
			return
		}

		if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", path, err)
			return
		}

		fmt.Println("processed", path)
	}
}
