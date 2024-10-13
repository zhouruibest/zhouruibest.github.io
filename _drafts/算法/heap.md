# 使用sort实现结构体排序

```go
import (
    "fmt"
    "sort"
)

type Person struct {
    Name string // 姓名
    Age  int    // 年纪
}

// 按照 Person.Age 从大到小排序
type PersonSlice []Person

func (a PersonSlice) Len() int { // 重写 Len() 方法
    return len(a)
}
func (a PersonSlice) Swap(i, j int) { // 重写 Swap() 方法
    a[i], a[j] = a[j], a[i]
}
func (a PersonSlice) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
    return a[i].Age < a[j].Age
}

people := []Person{...}
fmt.Println(people)
sort.Sort(PersonSlice(people)) // 按照 Age 的升序排序
sort.Sort(sort.Reverse(PersonSlice(people))) // 按照 Age 的降序排序

```

# 使用sort+container/heap实现int最小堆

```go
import (
    "container/heap"
    "fmt"
)

type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { // 注意，这里直接加到最后面
    *h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {   // 注意，这里弹出的是最后一个
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

h := &IntHeap{2, 1, 5, 100, 3, 6, 4, 5}
heap.Init(h)
heap.Push(h, 3)

for h.Len() > 0 {
    fmt.Printf("%d ", heap.Pop(h))  // 这里就可以按照有小到大弹出
}
```
# 使用sort+container/heap实现结构体最小堆

```go

import (
    "container/heap"
)

type stu struct {
    name string
    age  int
}
type Stu []stu

func (t *Stu) Len() int {
    return len(*t) //
}

func (t *Stu) Less(i, j int) bool {
    return (*t)[i].age < (*t)[j].age
}

func (t *Stu) Swap(i, j int) {
    (*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

func (t *Stu) Push(x interface{}) {
    *t = append(*t, x.(stu))
}

func (t *Stu) Pop() interface{} {
    n := len(*t)
    x := (*t)[n-1]
    *t = (*t)[:n-1]
    return x
}


student := &Stu{{"Amy", 21}, {"Dav", 15}, {"Spo", 22}, {"Reb", 11}}
heap.Init(student)
one := stu{"hund", 9}
heap.Push(student, one)
for student.Len() > 0 {
    fmt.Printf("%v\n", heap.Pop(student))
}
```

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