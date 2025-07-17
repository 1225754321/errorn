package main

const (
	PACKAGE_NAME = "github.com/1225754321/errorn/errorn"
	IDENT_NAME   = "errorn"

	CONVERT_FUNC_NAME = "Of"

	FLAG_FUNC_NAME_DEFAULT = "Wrap"
	FLAG_FUNC_NAME_PANIC   = "Panic"
)

type FuncFlag int

const (
	FLAG_DEFAULT FuncFlag = iota // 默认处理函数
	FLAG_PANIC                   // panic处理函数
)

var (
	// flag标识函数逻辑相关
	FLAG_FUNC_NAME_MAP = map[string]FuncFlag{
		FLAG_FUNC_NAME_DEFAULT: FLAG_DEFAULT,
		FLAG_FUNC_NAME_PANIC:   FLAG_PANIC,
	}
)

var (
	// 支持在go.mod中使用包名别名场景
	packageName = PACKAGE_NAME
)

func SetPackName(name string) {
	packageName = name
}
