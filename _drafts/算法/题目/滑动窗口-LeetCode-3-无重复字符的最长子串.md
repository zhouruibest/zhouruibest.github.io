```go
func lengthOfLongestSubstring(s string) int {
    ans := 0

    c := []byte(s)
    n := len(c)
    left, right := 0, 0
    counter := map[byte]int{}

    for right < n {
        counter[s[right]]++
        for counter[s[right]] > 1 && left <= right {
            counter[s[left]]--
            left++
        }
        ans = max(ans, right-left+1)
        right++
    }
    
    return ans
}

```