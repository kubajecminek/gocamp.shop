package template

func MulFunc(a, b int) int {
	return a * b
}

func AddFunc(nums ...int) int {
	res := 0
	for _, n := range nums {
		res += n
	}
	return res
}
