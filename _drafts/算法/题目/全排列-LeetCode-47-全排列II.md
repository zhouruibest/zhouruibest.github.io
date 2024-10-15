```go
import "sort"

func permuteUnique(nums []int) (ans [][]int) {
    n := len(nums)

    sort.Ints(nums)
    
    perm := []int{}
    vis := make([]bool, n)

    var dfs func(int)
    dfs = func(idx int) {
        if idx == n {
            ans = append(ans, append([]int(nil), perm...))
            return
        }

        for i, v := range nums {
            if vis[i] || i > 0 && !vis[i-1] && v == nums[i-1] {
                continue
            }
        
            perm = append(perm, v)
            vis[i] = true
            dfs(idx + 1)
            vis[i] = false
            perm = perm[:len(perm)-1]
        }
    }
    dfs(0)

    return
}
```