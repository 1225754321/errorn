# errorn - Go错误处理简化库

## 功能
- 链式错误处理语法
- 通过`go generate`自动转换代码
- 通过`errorn`命令直接转换代码
- 支持多种返回场景（普通返回、panic等）
- 错误包装和上下文添加

## 安装errorn工具
```bash
go install github.com/1225754321/errorn
```

## 使用
1. 导入包并添加生成指令：
```go
//go:generate errorn
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/1225754321/errorn/errorn"
)
```

2. 使用链式错误处理语法：
```go
func example() error {
    err := someFunc()
    errorn.Wrap(err, func(e error) error {
        return fmt.Errorf("wrapped: %w", e)
    })
    errorn.Panic(err) // 如果err不为nil则panic
}
```

3. 添加默认错误处理函数（可选）：
```go
ee.AddDefaultOfFunc(func(err error) error {
    return fmt.Errorf("default handler: %w", err)
})
```

4. 运行生成命令：
```bash
go generate ./...
# 或者
errorn -d .
```

## 示例
```go
//go:generate errorn
package main

import (
    ee "errorn/errorn"
    "errors"
    "fmt"
    "log"
)

func main() {
    ee.AddDefaultOfFunc(func(err error) error {
        return fmt.Errorf("default handler: %w", err)
    })
    
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    err := step1()
    ee.Wrap(err, func(e error) error {
        return fmt.Errorf("step1 failed: %w", e)
    })
    ee.Panic(err)

    err = step2()
    ee.Wrap(err)

    return nil
}

func step1() error {
    return errors.New("original error")
}

func step2() error {
    return errors.New("another error")
}
```

## 贡献
欢迎提交PR，请确保：
1. 代码通过测试
2. 遵循现有代码风格
3. 更新相关文档
