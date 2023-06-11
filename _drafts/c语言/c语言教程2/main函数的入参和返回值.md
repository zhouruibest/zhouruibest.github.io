# 返回值
main函数的返回值，用于说明程序的退出状态。如果返回0，则代表程序正常退出；返回其他数字的含义则由系统决定，通常，返回非零代表程序异常退出，即使程序运行结果正确也仍需修复代码。

# 入参

void main()是错误的。C++之父 Bjarne Stroustrup 在他的主页上的 FAQ 中明确地写着 The definition void main( ) {}is not and never has been C++, nor has it even been C.（ void main( )从来就不存在于C++ 或者 C ）。

在最新的 C11 标准中，只有以下两种定义方式是正确的：
```
int main( void ) 
int main( int argc, char *argv[])
``` 
main 函数的返回值类型必须是 int ，这样返回值才能传递给程序的激活者（如操作系统）。 如果 main 函数的最后没有写 return 语句的话，C99 规定编译器要自动在生成的目标文件中（如 exe 文件）加入return 0；

## 例如
``` c
#include<stdio.h>

int main(int argc, char *argv[]) //一般使用argc来统计参个数，argv来存储参数具体值, 以空格隔开
{
	printf("argc is %d \n", argc);
 
	int i;
 
	for (i = 0; i<argc; i++)
	{
		printf("arcv[%d] is %s\n", i, argv[i]);
	}
	
	return 0;
}
```