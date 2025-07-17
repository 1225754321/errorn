//go:generate errorn
package main

import (
	"errors"
	fmt "fmt"
	"log"

	ee "github.com/1225754321/errorn/errorn"
)

func main() {
	ee.AddDefaultOfFunc(func(err error) error {
		return fmt.Errorf("default handler: %w", err)
	})
	// 模拟一个返回错误的函数
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

	// 无处理函数的Wrap
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
