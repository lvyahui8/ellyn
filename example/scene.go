package main

import (
	"errors"
	"example/dao/model"
	"fmt"
	"runtime"
	"time"
)

func init() {
	fmt.Printf("init example main")
}

func Sum(a, b int) int {
	time.Sleep(10 * time.Millisecond)
	return a + b
}

func Callback(fn func()) {
	fn()
}

func Handle() {
	_, _, _ = WithUnusualParam(1, 2, "ccc", true, 0)
	NoName(1, "xx")
	KeywordNameParam("x", "x", "x")
}

func Trade() {
	a := 0
	for i := 0; i < 10; i++ {
		a += i
	}
	switch a {
	case 10:
		a -= 1
	case 20:
		a -= 2
	default:
		a -= 3
	}

	defer func() {
		_ = recover()
	}()

	defer Sum(1, 2)

	time.Sleep(5 * time.Millisecond)

	if a > 10 {
		a += 10
	} else if a < 5 {
		a += 5
	} else {
		a += 1
	}

	for {
		a--
		if a <= 10 {
			break
		}
	}

	Callback(func() {
		k := 0
		for i := 0; i < 10; i++ {
			k++
		}
		if k > 10 {
			k = 0
		}
	})

	go func() {
		Handle()
	}()

	go Handle()

	c := make(chan string, 100)
	c <- "hello"
	select {
	case msg := <-c:
		_ = fmt.Sprintf("say: %v", msg)
	}
	if 1 == 1 {
	}
	func() {
	}()
	// empty
	go func() {
		select {} // 阻塞
	}()
}

func N(n int) int {
	if n == 0 {
		return 1
	}
	time.Sleep(10 * time.Millisecond)
	return n * N(n-1)
}

func WithUnusualParam(a, b int,
	c string, _ bool, bool int) (x, y int,
	_ error) {
	time.Sleep(10 * time.Millisecond)
	if bool == 0 {
		return 0, 0, errors.New("test")
	}
	return 0, 0, nil
}

func KeywordNameParam(bool, int, a string) (_ bool, string, uint float32) {
	return false, 0.0, 0.0
}

func NoName(int,
	string) (byte, bool) {
	_ = fmt.Sprintf("x%d", runtime.NumCPU())
	return 0, false
}

func NotCollectVars(arr [10]int, user model.User) {
	_ = fmt.Sprintf("arr %v", arr)
}

func ReDefinePtr() model.UserPtr {
	return nil
}
