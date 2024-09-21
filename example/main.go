package main

import "github.com/lvyahui8/ellyn/ellyn_agent"

import "fmt"

func Sum(a, b int) int {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 0)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	return a + b
}

func Callback(fn func()) {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 1)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	fn()
}

func Handle() {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 2)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)

}

func Trade() {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 3)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	a := 0
	for i := 0; i < 10; i++ {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
		a += i
	}
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 2)
	switch a {
	case 10:
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 3)
		a -= 1
	case 20:
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 4)
		a -= 2
	default:
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 5)
		a -= 3
	}

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 6)
	defer func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 4)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		_ = recover()
	}()

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 7)
	defer Sum(1, 2)

	if a > 10 {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 8)
		a += 10
	} else {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 9)
		if a < 5 {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 10)
			a += 5
		} else {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 11)
			{
				a += 1
			}
		}
	}

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 12)
	for {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 13)
		a--
		if a <= 10 {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 14)
			break
		}
	}

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 15)
	Callback(func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 5)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		k := 0
		for i := 0; i < 10; i++ {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
			k++
		}
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 2)
		if k > 10 {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 3)
			k = 0
		}
	})

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 16)
	go func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 6)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		Handle()
	}()

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 17)
	go Handle()

	c := make(chan string, 100)
	c <- "hello"
	select {
	case msg := <-c:
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 18)
		_ = fmt.Sprintf("say: %v", msg)
	}
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 19)
	if 1 == 1 {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 20)
	}
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 21)
	func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 7)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	}()
	// empty
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 22)
	go func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 8)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		select {} // 阻塞
	}()
}

func main() {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 9)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	fmt.Println(Sum(1, 1))
	Trade()
}
