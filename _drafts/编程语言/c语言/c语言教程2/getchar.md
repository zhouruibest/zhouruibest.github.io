```c
#include <stdio.h>

int main ()
{
   char c;
 
   printf("请输入字符：");
   c = getchar();
 
   printf("输入的字符：");
   putchar(c);

   return(0);
}
```

scanf() 可输入不包含空格的字符串，不读取回车，空格和回车表示输入完毕。
getchar() 只能读取用户输入缓存区的一个字符，包括回车。原型 int getchar(void) 从标准输入 stdin 获取一个字符（一个无符号字符）。这等同于 getc 带有 stdin 作为参数。