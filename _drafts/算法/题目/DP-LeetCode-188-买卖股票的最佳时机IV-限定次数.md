给你一个整数数组 prices 和一个整数 k ，其中 prices[i] 是某支给定的股票在第 i 天的价格。

设计一个算法来计算你所能获取的最大利润。你最多可以完成 k 笔交易。也就是说，你最多可以买 k 次，卖 k 次。

注意：你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。

```go
import (
	"math"
)

// 注意，由题目得知，买一次并卖掉算完成一次交易。那么在卖掉的时候给交易次数递减

func maxProfit(k int, prices []int) int {
    n := len(prices)
    // 回溯法
    // var dfs func(i, times int, hold bool) int
    // dfs = func(i, times int, hold bool) int {
    //     if times < 0 {
    //         return math.MinInt
    //     }

    //     if i < 0 {
    //         if hold {
    //             return math.MinInt
    //         } else {
    //             return 0
    //         }
    //     }

    //     if hold {
    //         return max(dfs(i-1, times, true), dfs(i-1, times, false) - prices[i])
    //     }
    //     return max(dfs(i-1, times, false), dfs(i-1, times-1, true) + prices[i])
    // }

    // return dfs(n-1, k, false)


    // 回溯法+记忆话搜索
    // 使用一个多维数组或者多级map保存dfs的中间结果，过程略

    // 递推/动态规划
    // 翻译成递推 times ==> j, true==>1, false==>0
    // （1）f[i][j][0] = max(f[i-1][j][0], f[i-1][j][1]+prices[i])
    // （2）f[i][j][1] = max(f[i-1][j][1], f[i-1][j-1][0])-prices[i])
    // 因为i,j是可以等于0的，那么（2）中会出现下标为-1的情况
    // 因此需要在f和f[i]的前面插入一个状态，因此最终的递推公式如下
    // (1) f[*][0][*] = math.MinInt, 这是针对j的优化
    // (2) f[0][j][0] = 0            这是针对i的优化
    // (3) f[0][j][1] = math.MinInt  这是针对i的优化

    f := make([][][]int, n+1)
    for i:=0; i<n+1;i++ {
        f[i] = make([][]int, k+2)
        for j:=0; j<k+2; j++ {
            f[i][j] = make([]int, 2)
            f[i][j][0] = math.MinInt // 第0维表示不持有股票
            f[i][j][1] = math.MinInt // 第1维表示持有股票
        }
    }

    // 上面已经都初始化成了负无穷，这里初始化成0即可
    for j:=1; j<k+2; j++ {
        f[0][j][0] = 0
    }

    for i:=0; i<n; i++ {
        for j:=1; j<k+2; j++ {
            f[i+1][j][0] = max(f[i][j][0], f[i][j-1][1]+prices[i])
            f[i+1][j][1] = max(f[i][j][1], f[i][j][0]-prices[i])
        }
    }

    return f[n][k+1][0]
    

    // 递推/动态规划 优化内存

}
```