package tools

// NumberSqrt MySqrt 使用牛顿法求平方根
func NumberSqrt(x int) int {
	res := x
	//牛顿法求平方根
	for res*res > x {
		res = (res + x/res) / 2
	}
	return res
}