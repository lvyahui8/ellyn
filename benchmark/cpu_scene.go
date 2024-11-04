package benchmark

import "benchmark/ellyn_agent"

func partition(arr []int, low, high int) ([]int, int) {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 0, nil)
	defer ellyn_agent.Agent.Pop(_ellynCtx, nil)
	ellyn_agent.Agent.SetBlock(_ellynCtx, 0, 0)
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		ellyn_agent.Agent.SetBlock(_ellynCtx, 1, 2)
		if arr[j] < pivot {
			ellyn_agent.Agent.SetBlock(_ellynCtx, 2, 3)
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	ellyn_agent.Agent.SetBlock(_ellynCtx, 3, 1)
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func quickSort(arr []int, low, high int) []int {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 1, nil)
	defer ellyn_agent.Agent.Pop(_ellynCtx, nil)
	ellyn_agent.Agent.SetBlock(_ellynCtx, 0, 4)
	if low < high {
		ellyn_agent.Agent.SetBlock(_ellynCtx, 1, 6)
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	ellyn_agent.Agent.SetBlock(_ellynCtx, 2, 5)
	return arr
}

func quickSortStart(arr []int) []int {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 2, nil)
	defer ellyn_agent.Agent.Pop(_ellynCtx, nil)
	ellyn_agent.Agent.SetBlock(_ellynCtx, 0, 7)
	return quickSort(arr, 0, len(arr)-1)
}
