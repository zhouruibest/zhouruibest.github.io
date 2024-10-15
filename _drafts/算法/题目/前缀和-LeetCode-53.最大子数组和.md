```go
import "math"

func maxSubArray(nums []int) int {
    ans := math.MinInt
    minPreSum := 0
    preSum := 0
    for _, x := range nums {
        preSum += x // 当前的前缀和
        ans = max(ans, preSum-minPreSum)   // 减去前缀和的最小值
        minPreSum = min(minPreSum, preSum) // 维护前缀和的最小值
    }
    return ans
}
```