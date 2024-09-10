因为 if-else 语句的 else 部分是可选的，所以在嵌套的 if 语句中省略它的 else 部 分将导致歧义。解决的方法是将每个 else 与最近的前一个没有 else 配对的 if 进行匹配。

例如，在下列语句中:

```c
if (n > 0)
    if (a > b)
        z = a;
else
    z = b;
```

此时 **else 部分与内层的 if 匹配**.