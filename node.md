# golang error简易操作库

我有个想法,我想要简便golang的error操作:
例如:
原来:
func xxx() error {
    err := xx()
    if err != nil {
        return err
    }
}
简化后:
//go:generate errorN
package xxxxx
import (
    errorN "github.com/xxx/errorN"
)
func xxx() error {
    err := xx()
    errorN.Of(err).Chain(func(err)err)
}

写完上面的代码后
执行go generate ./...

我会将上面的代码进行替换为
//go:generate errorN
package xxxxx
import (
    errorN "github.com/xxx/errorN"
)
func xxx() error {
    err := xx()
    if err != nil {
        return errorN.From(err).Chain(xxxx)
    }
}