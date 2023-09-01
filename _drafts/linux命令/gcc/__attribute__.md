# C语言的__attribute__机制
使GUN的扩展，不是C标准
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
```c
__attribute__((noreturn))
void test()
{
    exit(0); // 注意，这里不能return；
}
```

## 用法3
设置整形的长度
```c
typedef unsigned int myint __attribute((mode(HI)));
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