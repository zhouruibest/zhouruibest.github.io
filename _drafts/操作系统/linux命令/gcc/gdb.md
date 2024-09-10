# 准备

为了能让GDB调试，可执行文件需要附加额外的调试信息，因此在编译的时候需要给GCC传递-g参数如
```sh
gcc -g demo-gdb.c -o a.out
```

# 简单执行

```sh
gdb a.out
```

# 常用指令

## x, examine

