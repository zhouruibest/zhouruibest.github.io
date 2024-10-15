```go
func combine(n int, k int) [][]int {
    ans := make([][]int, 0)
    path := []int{}

    var dfs func(i int)
    dfs = func(i int) {
        // 到达结尾
        if i == n {
            if len(path) == k {
                path2 := make([]int, len(path))
                copy(path2, path)
                ans = append(ans, path2)
            }
            return
        }

        // 没有到达结尾
        if len(path) == k {
            dfs(n)
        } else if len(path) + (n-i) < k {
            //
        } else if len(path) + (n-i) == k{
            for j:=i; j<n; j++ {
                path = append(path, j+1)
            }
            dfs(n)
        } else {
            pathLen := len(path)
            path = append(path, i+1)
            dfs(i+1)
            path = path[:pathLen]
            dfs(i+1)
        }
    }
    dfs(0)
    return ans
}
```