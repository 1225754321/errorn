package errorn

// Wrap 是开发阶段的占位函数，用于标记错误处理点
// 代码生成器会将其替换为 errorN.Of 调用
// 实际运行时不会执行任何操作
func Wrap(err error, handlers ...func(error) error) {
	// 空实现
}

func Panic(err error, handlers ...func(error) error) {
	// 空实现
}
