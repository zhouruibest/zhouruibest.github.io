list定义了一个双向链表

```go
package main

import (
    "container/list"
    "fmt"
)

func main() {
    l := list.New() // 创建双向链表

    l.PushBack("a") // 尾部追加元素
    printList(l) // a

    l.PushBack("b") // 尾部追加元素
    printList(l) // a b

    l.PushFront("c")
    printList(l) // c a b

    fmt.Println(l.Front().Value) // c 打印头部元素
    fmt.Println(l.Back().Value)  // b
    fmt.Println(l.Len())         // 3 链表长度

    l.MoveToBack(l.Front())
    printList(l) // a b c

    l.MoveToFront(l.Back())
    printList(l) // c a b

    l.Remove(l.Back())
    printList(l) // c a
}

func printList(l *list.List) {
    for e := l.Front(); e != nil; e = e.Next() {  // 遍历
        fmt.Print(e.Value, " ")
    }
    fmt.Println()
}
```