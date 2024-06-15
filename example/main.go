package main

import "fmt"

func Sum(a, b int) int {
	return a + b
}

func Callback(fn func()) {
	fn()
}

func Handle() {

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
}

func main() {
	fmt.Println(Sum(1, 1))
}
