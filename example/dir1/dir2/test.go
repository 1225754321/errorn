package main

import (
	"fmt"

	"github.com/1225754321/errorn/errorn"
)

// 零返回值函数
func zeroReturn() {
	err := fmt.Errorf("error") // 多error返回值 testnode
	errorn.Wrap(err)
	fmt.Println("test")
}

// 命名返回值函数
func namedReturn() (a int, b string, err error) {
	err = fmt.Errorf("error")
	errorn.Wrap(err) // 多error返回值 test
	return
}

// 多error返回值
func multiReturn() (int, string, error) {
	err := fmt.Errorf("error")
	errorn.Wrap(err) // 多error返回值
	return 0, "", nil
}

// 单error返回值
func singleReturn() error {
	err := fmt.Errorf("error")
	errorn.Wrap(err)
	return nil
}
