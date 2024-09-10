# C语言的__attribute__机制
是GUN的扩展，不是C标准
## 用法1
以下函数中，虽然before和after都没有被main调用， 但是他们都在main之前、之后调用了一下
```c
#include <stdio.h>
#include <stdlib.h>

__attribute__((constructor))
void before(void)
{
    printf("this is before...\n");
}

__attribute__((destructor))
void after(void)
{
    printf("this is after...\n");
}

int main(){
    printf("this is main...\n");
}

输出

this is before...
this is main...
this is after...

```
## 用法2
表明函数调用完成后不返回主调函数（比如里面有exit()，不是说没有返回值，而是指，不把控制权交给主调函数）, 注意，这与void返回类型不同。
目的1 用户不要随便调用这个函数
目的2 编译器可以优化一些代码

```c
// 用在函数定义
__attribute__((__noreturn__))
void test()
{
    exit(0); // 注意，这里不能return；
}

// 还可以用在函数声明
void usageErr(const char *format, ...) __attribute__((__noreturn__))
// 定义时就不用了
void
usageErr(const char *format, ...)
{
    ...
}

```

## 用法3
设置整形的长度
```c
typedef unsigned int myint __attribute__((mode(HI)));
```

## 用法4
放弃字节对齐

```c
struct Test {
    char ch;
    int id;
    double len;
}__attribute__((packed)); // 放弃对齐


struct Test2 {
    char ch;
    int id;
    double len;
}__attribute__((alifned(32))) // 制定对其宽度

```

## 用法5
静态函数如果没有调用就会警告，可以用
```c
__attribute__((unused))
static void pf(){

}
```