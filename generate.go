package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

// 创建errorN.Of调用表达式
func createOfCall(identName string, errIdent *ast.Ident, handlers []ast.Expr) *ast.CallExpr {
	args := []ast.Expr{errIdent}
	args = append(args, handlers...)
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: identName},
			Sel: &ast.Ident{Name: CONVERT_FUNC_NAME},
		},
		Args: args,
	}
}

// 获取函数返回类型信息
type returnInfo struct {
	types []ast.Expr
	names []string // 命名返回值名称
}

func getReturnInfo(fn ast.Node) *returnInfo {
	info := &returnInfo{}
	switch node := fn.(type) {
	case *ast.FuncDecl:
		if node.Type.Results != nil {
			for _, field := range node.Type.Results.List {
				info.types = append(info.types, field.Type)
				if len(field.Names) > 0 {
					info.names = append(info.names, field.Names[0].Name)
				} else {
					info.names = append(info.names, "")
				}
			}
		}
	case *ast.FuncLit:
		if node.Type.Results != nil {
			for _, field := range node.Type.Results.List {
				info.types = append(info.types, field.Type)
				if len(field.Names) > 0 {
					info.names = append(info.names, field.Names[0].Name)
				} else {
					info.names = append(info.names, "")
				}
			}
		}
	}
	return info
}

// 检查是否有命名返回值
func hasNamedReturns(info *returnInfo) bool {
	for _, name := range info.names {
		if name != "" {
			return true
		}
	}
	return false
}

// 验证返回类型
func validateReturnTypes(returnTypes []ast.Expr) error {
	if len(returnTypes) == 0 {
		return nil // 零返回值函数有效
	}

	// 检查最后一个返回类型是否为error
	lastType := returnTypes[len(returnTypes)-1]
	if !isErrorType(lastType) {
		return fmt.Errorf("last return type must be error, got %s", exprToString(lastType))
	}
	return nil
}

// 检查是否为error类型
func isErrorType(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "error"
	}
	return false
}

// 表达式转字符串
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}

// 生成替换语句
func generateReplacement(flag FuncFlag, identName string, errIdent *ast.Ident, handlers []ast.Expr, returnInfo *returnInfo) ast.Stmt {
	returnTypes := returnInfo.types
	if len(returnTypes) == 0 || flag == FLAG_PANIC {
		// 零返回值函数 -> panic
		return &ast.IfStmt{
			If: token.NoPos,
			Cond: &ast.BinaryExpr{
				X:  errIdent,
				Op: token.NEQ,
				Y:  &ast.Ident{Name: "nil"},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun:  &ast.Ident{Name: "panic"},
							Args: []ast.Expr{createOfCall(identName, errIdent, handlers)},
						},
					},
				},
			},
		}
	}

	// 处理命名返回值场景
	if hasNamedReturns(returnInfo) {
		// 创建命名返回值列表
		var results []ast.Expr
		for i := 0; i < len(returnTypes)-1; i++ {
			if returnInfo.names[i] != "" {
				results = append(results, &ast.Ident{Name: returnInfo.names[i]})
			} else {
				results = append(results, zeroValue(returnTypes[i]))
			}
		}
		// 最后一个返回值是error
		results = append(results, createOfCall(identName, errIdent, handlers))

		return &ast.IfStmt{
			If: token.NoPos,
			Cond: &ast.BinaryExpr{
				X:  errIdent,
				Op: token.NEQ,
				Y:  &ast.Ident{Name: "nil"},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{Results: results},
				},
			},
		}
	}

	// 非命名返回值场景
	results := make([]ast.Expr, len(returnTypes))
	for i := 0; i < len(returnTypes)-1; i++ {
		results[i] = zeroValue(returnTypes[i])
	}
	results[len(returnTypes)-1] = createOfCall(identName, errIdent, handlers)

	return &ast.IfStmt{
		If: token.NoPos,
		Cond: &ast.BinaryExpr{
			X:  errIdent,
			Op: token.NEQ,
			Y:  &ast.Ident{Name: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{Results: results},
			},
		},
	}
}

// 生成类型零值
func zeroValue(expr ast.Expr) ast.Expr {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "bool":
			return &ast.Ident{Name: "false"}
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64", "complex64", "complex128":
			return &ast.BasicLit{Kind: token.INT, Value: "0"}
		case "string":
			return &ast.BasicLit{Kind: token.STRING, Value: `""`}
		case "error":
			return &ast.Ident{Name: "nil"}
		}
	}
	return &ast.Ident{Name: "nil"} // 默认返回nil
}

func generateDefault(flag FuncFlag, identName string, errIdent *ast.Ident, handlers []ast.Expr) ast.Stmt {
	switch flag {
	case FLAG_PANIC:
		return &ast.IfStmt{
			If: token.NoPos,
			Cond: &ast.BinaryExpr{
				X:  errIdent,
				Op: token.NEQ,
				Y:  &ast.Ident{Name: "nil"},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun:  &ast.Ident{Name: "panic"},
							Args: []ast.Expr{createOfCall(identName, errIdent, handlers)},
						},
					},
				},
			},
		}
	}
	return &ast.IfStmt{
		If: token.NoPos,
		Cond: &ast.BinaryExpr{
			X:  errIdent,
			Op: token.NEQ,
			Y:  &ast.Ident{Name: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						createOfCall(identName, errIdent, handlers),
					},
				},
			},
		},
	}
}
