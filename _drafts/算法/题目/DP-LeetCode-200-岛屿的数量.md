```go

func numIslands(grid [][]byte) int {
    ans := 0

    for r:=0; r<len(grid); r++ {
        for c:=0; c<len(grid[0]); c++ {
            if grid[r][c] == '1' {
                mark(grid, r, c)
                ans++
            }
        }
    }

    return ans
}

func mark(grid [][]byte, r, c int) {
   if !inArea(grid, r, c) { // 深度优先遍历， 来到了边界
        return
   }

   if grid[r][c] != '1' { // 0海洋； 1陆地； 2陆地但是遍历过了
        return
   }

   grid[r][c] = '2'
   mark(grid, r-1, c)
   mark(grid, r+1, c)
   mark(grid, r, c-1)
   mark(grid, r, c+1)
}

func inArea(grid [][]byte, r, c int) bool {
    return r >=0 && r<len(grid) && c>=0 && c<len(grid[0])    
}
```