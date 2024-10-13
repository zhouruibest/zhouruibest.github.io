# 求解最大公约数

```go
func gcd(x, y int) int {
	tmp := x % y
	if tmp > 0 {
		return gcd(y, tmp)
	}
	return y
}
```

# 求解最小公倍数

```go
func lcm(x, y int) int {
	return x * y / gcd(x, y)
}
```
