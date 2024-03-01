mmap函数是一个比较神奇的函数，它可以把文件映射到进程的虚拟内存空间。通过对这段内存的读取和修改，可以实现对文件的读取和修改，而不需要用read和write函数。如下图所示，为mmap实现原理的示意图。

2、下面我们来看一下mmap函数的原型

void *mmap(void *addr, size_t len, int prot, int flags, int fd, off_t offset);
在这个函数原型中：

参数addr：指定映射的起始地址，通常设为NULL，由内核来分配

参数length：代表将文件中映射到内存的部分的长度。

参数prot：映射区域的保护方式。可以为以下几种方式的组合：

PROT_EXEC 映射区域可被执行
PROT_READ 映射区域可被读取
PROT_WRITE 映射区域可被写入
PROT_NONE 映射区域不能存取

参数flags：映射区的特性标志位，常用的两个选项是：

MAP_SHARD：写入映射区的数据会复制回文件，且运行其他映射文件的进程共享

MAP_PRIVATE：对映射区的写入操作会产生一个映射区的复制，对此区域的修改不会写会原文件

参数fd：要映射到内存中的文件描述符，有open函数打开文件时返回的值。

参数offset：文件映射的偏移量，通常设置为0，代表从文件最前方开始对应，offset必须是分页大小的整数倍。

函数返回值：实际分配的内存的起始地址。

3、munmap函数

        与mmap函数成对使用的是munmap函数，它是用来解除映射的函数，原型如下：

int munmap(void *start, size_t length)
在这个函数中，

参数start：映射的起始地址

参数length：文件中映射到内存的部分的长度

返回值：解除成功返回０，失败返回-1。


```c
	//打开文件
	fd = open("testdata",O_RDWR);
	//创建mmap
	start = (char *)mmap(NULL,128,PROT_READ|PROT_WRITE,MAP_SHARED,fd,0);
	//读取文件	
	strcpy(buf,start);
	printf("%s\n",buf);
	//写入文件
	strcpy(start,"Write to file!\n");
 
	munmap(start,128);
	close(fd);
    ```