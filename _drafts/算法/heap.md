# 使用sort.IntSlice实现最小堆

```go

import (
	"container/heap"
	"sort"
)

var factors = []int{2, 3, 5}

type hp struct{sort.IntSlice}
func (h *hp) Push(v interface{}) {h.IntSlice = append(h.IntSlice, v.(int))} // 注意，这里是直接追加到最后一个的
func (h *hp) Pop() interface{} {a := h.IntSlice; v := a[len(a)-1]; h.IntSlice = a[:len(a)-1]; return v} // 注意，这里Pop的是最后一个，返回的也是最后一个

// 初始化
h := &hp{sort.IntSlice{1}}

// 弹出元素
x := heap.Pop(h).(int)

// 加入元素
heap.Push(h, next)
```