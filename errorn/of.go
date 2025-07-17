package errorn

// Of 处理错误并应用处理函数链
// 如果err为nil则返回nil
// 如果有处理函数，按顺序执行：handler1(handler2(...handlerN(err)))
// 如果没有处理函数，直接返回原始错误
func Of(err error, handlers ...func(error) error) error {
	if err == nil {
		return nil
	}

	// 没有处理函数时直接返回原始错误
	if len(handlers) == 0 {
		return err
	}

	// 按顺序执行处理函数
	current := err
	for _, handler := range handlers {
		current = handler(current)
	}

	for _, df := range defaultOfFunc {
		if df == nil {
			continue
		}
		current = df(current)
	}

	return current
}

var (
	defaultOfFunc = make([]func(error) error, 0, 5)
)

func AddDefaultOfFunc(fs ...func(error) error) {
	defaultOfFunc = append(defaultOfFunc, fs...)
}

func UpdateDefaultOfFunc(fs []func(error) error) {
	defaultOfFunc = fs
}
