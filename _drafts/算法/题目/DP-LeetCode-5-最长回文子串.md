```go
func longestPalindrome(s string) string {
    // 任意一个字符都是长度为1的回文子串
    // 如果他的左右两个相同，那么构成长度为3的回文子串
    // 以此类推
    // dp

    n := len(s)
    data := []byte(s)

    dp := make([][]bool, n)
    for i := range dp {
        dp[i] = make([]bool, n)
    }

    maxLen := 1
    left, right := 0, 0 // 记录结果

    // 任意一个字符都是长度为1的回文子串
    for i := 0; i<n; i++ {
        dp[i][i] = true
    }
    // 处理相邻
    for i:=0; i<n-1;i++ {
        if s[i] == s[i+1] {
            dp[i][i+1] = true
            if dp[i][i+1] {
                maxLen = 2
                left = i
                right = i + 1
            }
        }
    }

    // 注意遍历次序
    for j := 2; j<n; j++ {
        for i:= 0; i+2<=j; i++ {
            dp[i][j] = s[i] == s[j] && dp[i+1][j-1]
            
            if dp[i][j] {
                length := j - i + 1
                if length > maxLen {
                    left = i
                    right = j
                    maxLen = length
                }
            }
        }
    }

    return string(data[left:right+1])

}
```