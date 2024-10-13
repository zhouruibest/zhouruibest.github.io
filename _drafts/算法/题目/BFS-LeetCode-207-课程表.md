```go
// 广度优先遍历时候一定要注意：到达点B的箭头线有几条？
// 如果从任意一个箭头线可以到达，那就很简单，cur->next,一轮一轮的迭代即可
// 如果必须同时从两个箭头到达（也就是这个节点有两个依赖），那就要考虑“入度”+“DAG”，是经典的图遍历问题
func canFinish(numCourses int, prerequisites [][]int) bool {
    var (
        edges = make([][]int, numCourses) // 有向连接线，numCourses[x][]表示从x点出发可以到达的集合
        indeg = make([]int, numCourses)   // 每一个点的入度
        result []int                      // 当前节点是否到达
    )

    // 初始化
    for _, info := range prerequisites {
        edges[info[1]] = append(edges[info[1]], info[0])
        indeg[info[0]]++
    }

    // 找出入度为0的节点
    q := []int{}
    for i := 0; i < numCourses; i++ {
        if indeg[i] == 0 {
            q = append(q, i)
        }
    }

    for len(q) > 0 {
        u := q[0]
        q = q[1:]
        result = append(result, u)
        for _, v := range edges[u] { // v表示当前节点u可以到达的节点v
            indeg[v]-- //节点v的入度-1
            if indeg[v] == 0 { // 如果节点v的入度减少至0，那么v进入队列
                q = append(q, v)
            }
        }
    }

    return len(result) == numCourses
}
```