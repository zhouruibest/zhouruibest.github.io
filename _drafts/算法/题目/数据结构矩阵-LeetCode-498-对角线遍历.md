```go
func findDiagonalOrder(mat [][]int) []int {
    m, n := len(mat), len(mat[0])

    ans := make([]int, 0, m*n)


    for i := 0; i < m+n-1; i++ {  // 一共要遍历m+n-1趟
        if i%2 == 1 { // 奇数 趟，从右上到左下
            x := max(i-n+1, 0)    
            y := min(i, n-1)
            for x < m && y >= 0 {
                ans = append(ans, mat[x][y])
                x++
                y--
            }
        } else {
            x := min(i, m-1)
            y := max(i-m+1, 0)
            for x >= 0 && y < n {
                ans = append(ans, mat[x][y])
                x--
                y++
            }
        }
    }

    return ans
}
```