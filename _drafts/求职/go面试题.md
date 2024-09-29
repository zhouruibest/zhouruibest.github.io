# context包用来干什么， BackGround 和 TODO两种类型的context有啥区别

用途：
基于 context 来实现和搭建了各类 goroutine 控制的，并且与 select-case联合，就可以实现进行上下文的截止时间、信号控制、信息传递等跨 goroutine 的操作

在 Golang 的 context 包中，context.TODO() 和 context.Background() 是两个常用的函数，它们有一些区别。

1. context.TODO()：
context.TODO() 返回一个空的 Context，用于表示在当前情况下无法确定应该使用哪种 Context。
当你在编写代码时，不确定应该使用哪种特定的 Context 类型时，可以使用 context.TODO()。
使用 context.TODO() 应该是一个临时的解决方案，你应该在明确需要特定类型的 Context 时尽早替换为适当的 Context。

2. context.Background()：
context.Background() 返回一个空的 Context，通常用作整个请求的顶级 Context。
它是一个最基本的 Context，没有任何附加值或取消信号。
当你无需传递具体的 Context 时，可以使用 context.Background()。

这两个函数的主要区别在于 context.TODO() 表示**暂时的、未确定的** Context，而 context.Background() 表示**顶级的、无附加值的基本 Context**。在实际使用中，context.TODO() 应该尽可能被替换为适当的 Context 类型，而 context.Background() 可以作为通用的默认选择。

需要根据具体的情况选择使用合适的 Context 类型，以满足代码逻辑和需求。

