# 本地化

很多人认为C语言只支持ascii码，这是误解，你用printf(“这是中文”)同样可以输出中文，用fputs(“C语言本地化”)也可以向文件写入中文。那么C语言默认使用什么编码呢？这个问题不是那么简单，因为C标准并未规定使用什么编码，实际编码与操作系统、所在区域、编译器有很大的关系。另外可能你会发现虽然可以输出中文字符串，但并不能用char c=’中’来声明一个中文字符，如果用strlen(“C语言”)来求字符串长度，给出的值可能是5，因为编码是变长的，默认情况下C并不能很好的支持中文。如何让C语言很好的支持中文呢？这需要我们从基本编码概念讲起

# 字符编码

计算机刚刚发明时，只支持ascii码，也就是说只支持英文，随着计算机在全球兴起，各国创建了属于自己的编码来显示本国文字，中文首先使用GB2132编码（收录了6763个汉字）、GBK（21003个汉字）、GB18030（收录了27533个汉字），**这些中文编码本身兼容ascii，并采用变长方式记录**，英文使用一个字节，常用汉字使用2个字节，罕见字使用四个字节。后来随着全球文化不断交流，人们迫切需要一种全球统一的编码能够统一世界各地字符，再也不会因为地域不同而出现乱码，这时**Unicode字符集**就诞生了，也称为统一码，万国码。新出来的操作系统其内核本身就支持Unicode。考虑到空间和性能，Unicode提供了3种编码方案：
- utf-8 变长编码方案，使用1-6个字节来储存
- utf-32 定长编码方案，始终使用4个字节来储存
- utf-16 介于变长和定长之间的平衡方案，使用2个或4个字节来储存
utf-8由于是变长方案，类似GB2132和GBK量体裁衣，最节省空间，但**要通过第一个字节决定采用几个字节储存**，编码最复杂，且由于变长要定位文字，就得从第一个字符开始计算，性能最低。utf-32由于是定长方案，字节数固定因此无需解码，性能最高但最浪费空间。utf-16是个怪胎，它将常用的字放在编号0 ~ FFFF之间，不用进行编码转换，对于不常用字的都放在10000~10FFFF编号之后，因此自然的解决变长的问题。**对于这3种编码，只有utf-8兼容ascii，utf-32和utf-16都不支持单字节**，由于utf-8最省流量，兼容性好，后来解码性能也得到了很大改善，同时新出来的硬件也越来越强，性能已不成问题，因此很多纯文本、源代码、网页、配置文件等都采用utf-8编码，从而代替了原来简陋的ascii。再来看看utf-16，对于常见字2个字节已经完全够用，很少会用到4个字节，因此通常也将utf-16划分为定长，一些操作系统和代码编译器直接不支持4字节的utf-16。**Unicode还分为大端和小端**，大端就是将高位的字节放在低地址表示，后缀为BE；小端就是将高位的字节放在高地址表示，后缀为LE，没有指定后缀，即不知道其是大小端，所以其开始的两个字节表示该字节数组是大端还是小端，FE FF表示大端，FF FE表示小端。**Windows内核使用utf-16，linux，mac，ios内核使用的是utf-8**，我们就不去争论谁好谁坏了。另外虽然windows内核为utf-16，但为了更好的本地化，控制面板提供了区域选项，如果设置为简体就是GBK编码，在win10中，控制台和记事本默认编码为gbk，其它第三方软件就不好说了，它们默认编码各不相同。

了解编码后下面来说说BOM，一个文本文件，可以是纯文本、网页、源码等，那么打开它的程序如何知道它采用什么编码呢？为了**说明一个文件采用的是什么编码，在文件最开始的部分，可以有BOM**，比如0xFE 0xFF表示UTF-16BE，0xFF 0xFE 0x00 0x00表示UTF-32LE。UTF-8原本是不需要BOM的，因为其自我同步的特性，但是为了明确说明这是UTF-8而不是让文本编辑器去猜，也可以加上UTF-8的BOM：0xEF 0xBB 0xBF

# C语言字符编码

C语言源代码采用什么格式取决于IDE环境，通常是utf-8或ANSI，什么是**ANSI编码**呢？相比unicode它是采取另一种思路，严格来说ANSI不是一种编码，而是一种替代方案，力求找到显示内容的最低编码需求，如果内容只有英文字符就使用ascii，如果发现汉字就替换成本地的GBK编码，如果发现既有汉字又有日语又有韩语是否会自动选择unicode(这个没有试过)。前面讲过C语言使用的编码和操作系统、区域选择、编译器都有关系，但有一个现象，通常源码采用什么格式的编码，运行时就使用这样的编码，因此我们可以通过源代码先看看这个IDE会使用什么编码。在Dev C++中测试的结果是，如果只有英文源码默认采用utf8，现在的编辑器很少使用纯ascii了，如果发现源码里面有汉字，则将编码改为ANSI，由于windows控制台默认也使用ANSI，因此可以显示标准输入输出中的汉字。

如果只需要向控制台输出一段字符串或者向文件中写入一段中文，那么使用标准输入输出函数即可，**因为printf()和puts()可以识别字符串使用的编码**，但如果要操作单个字符就不行了，因为GBA和GBK都是变长的，一个英文字母使用一个字节，而一个中文使用2-4个字节，用char c=’中’是不行的，因为char只能是一个字节。**再来看看字符串，虽然可以放入中文，但却不能访问数组元素，因为数组元素也只能是一个char，数组元素个数和字符串个数是不对应的**，这是历史遗留的问题，无法更改，否则新标准就不能兼容之前的代码了。要处理中文字符，只能另辟途径，一是将编码从变长改为定长，二是使用新的字符类型来处理定长编码，对于变长和定长，在计算机行业种还有一个术语称为*窄字符和宽字符*，从前面知识可以得知，能够显示各个国家的语言并且采用定长的只有utf-16和utf-32了，C语言为宽字符提供新的类型wchar，这个类型由wchar库提供，**导入wchar.h后就可以使用宽字符了，宽字符使用utf-16还是utf-32由编译器决定**，windows中宽字符使用utf-16，linux则使用utf-32，我们可以通过以下代码进行测试，如下：

```c
#include <stdio.h>
#include <wchar.h>

int main()
{
	wchar_t wc1=L'a';
	wchar_t wc2=L'中';
	printf("%d,%d\n",sizeof(wc1),sizeof(wc2));
	printf("%x,%x\n",wc1,wc2);
	return 0;
}
```

上面代码使用wchar_t声明一个宽字符，宽字符前面要添加L，L是有个宏，它将后面的内容转为宽字符或宽字符串，这里使用sizeof()检测的结果是无论中文还是英文都是2字节，说明编码使用的是utf-16，后面一个printf()将编码以十六进制输出。至于utf-16和utf-32哪个好我们也不去争论了，巨头们都很任性且互不买账。

如果要显示一个宽字节符，需要换成putwchar()和wprintf()函数，如

```c
#include <stdio.h>
#include <wchar.h>
#include <locale.h>

int main()
{
	setlocale(LC_ALL, "");
	wchar_t wc1=L'a';
	wchar_t wc2=L'中';
	putwchar(wc1);
	putwchar(wc2);
	wprintf(L"%c%c",wc1,wc2);
	return 0;
}
```

这两个方法需要先调用`setlocale()`函数进行初始化，设置Unicode区域，setlocale()的格式为：
```c
char* setlocale (int category, const char* locale)
```
category是类型，表示区域编码影响到的类型，类型有时钟、货币、字符排序等，通常设置为常数LC_ALL表示影响所有类型。locale表示区域，windows和linux表示区域的方式各不相同，例如windows表示简体中文用”chs”，linux用”zh_CN”，但有3个区域总是相同的：”C”表示中立地域，不表示任何一个地区，只对小数点进行了设置，是默认值；“”表示本地区域；NULL表示不指定区域仅仅返回区域信息，可以通过`puts(setlocale(NULL,""))`输出本地区域信息，中文简体为Chinese (Simplified)_China.936。如果不设置区域，默认区域为”C”，这个区域只能显示英文，无法显示任何中文字符，实际测试结果也是如此。由于setlocale()包含在locale.h中，因此要先导入locale库，下面我们再来看看宽字符串的定义和输出：

```c
#include <stdio.h>
#include <wchar.h>
#include <locale.h>

int main()
{
        char *local = setlocale(LC_ALL, "");
        //printf("默认地域设置  %s \n", local);   // 注意1. printf和wprintf是不能混用的！！！！

        local = setlocale(LC_ALL, "zh_CN.UTF-8");
        //printf("修改后地域设置  %s \n", local);

        wchar_t wstr[] = L"宽字符串";
        wprintf(L"%ls,\n\n", wstr);

        //printf("长度 %d\n",wcslen(wstr));
        wchar_t str[] = L"你好，世界！";
        wprintf(L"xxxxxx %ls\n", str);
        return 0;
}

```
第二个例子
```c
#include <stdio.h>
#include <wchar.h>
#include <locale.h>

int main()
{
	setlocale(LC_ALL, "");
	wchar_t wc1=L'a';
	wchar_t wc2=L'中';
	putwchar(wc1);
	putwchar(wc2);
	wprintf(L"%c%c",wc1,wc2);
	return 0;
}
```

声明宽字符串后应该使用wcslen()函数获取字符串的长度，使用下标和指针运算都能正确定位，这就是采用定长的原因。除此之外，字符串的复制、连接等都有配套的宽字符串操作方法，相关方法可以查询C语言函数库。需要清楚的是，传统的字符串和相应的方法现在归于窄字符范畴，使用L前缀的字符串和带w的函数属于宽字符串范畴，它们的类型和处理方式都不相同，**不能混用!**, **不能混用!**, **不能混用!**。


# 宽字符的实现原理

当我们用大写的L标记一个宽字符串时，这个字符串编码会以utf-16或utf-32储存。然而，输出字符串时控制台或文本不一定是utf-16和utf-32，因此在输出时会转化为窄字符使用的编码，转化的依据是setlocale()中对宽字符编码的区域的设定，例如setlocale()将区域设定为本地，本地编码为简体中文，那么在windows中运行时会将zh-cn的unicode转为gbk，在linux中运行会转化为utf-8，实际上无论将控制台设置为什么编码，宽字符都能正常显示。这种编码的转换被写入到stout中，因此调用setlocale()后除了影响宽字符还影响窄字符，GCC就是一个例子，GCC对宽字符的支持很不友好，如果代码为ANSI或GBK，则以L开头的宽字符或宽字符串不能通过编译，**因为linux下的GCC在默认情况下总是假设源代码的编码是utf-8，只有将源码格式设置为utf-8才能通过编译**。：

如果将源码改为utf-8，则编译后输出的编码也为utf-8，还需要将控制台编码设置为utf-8才能显示窄字符，宽字符依然需要调用setlocale()设置区域，然而更改区域又会影响窄字符的显示，要同时显示窄字符和宽字符只能每次显示之前先设置区域，如下：

```c
#include <stdio.h>
#include<Windows.h>
#include<locale.h>

int main()
{
	system("chcp 936");//将控制台编码设置为gbk
	puts("乱码");
	system("pause");
	system("chcp 65001");//将控制台编码设置为utf8
	printf("%s\n", "显示窄字符utf8");
	setlocale(LC_ALL, "");
	printf("%ls\n",L"显示宽字符utf16");
	
	setlocale(LC_ALL, "C");
	puts("还原区域显示窄字符");
	return 0;
}

```

在windows控制台中执行chcp可以设置控制台字符编码，这里使用system(“chcp 936”)来调用这个命令将控制台编码设置为gbk，由于gcc输出编码为utf-8，因此在gbk下呈现乱码，接着调用system(“chcp 65001”)将控制台编码设置为utf8发现能够正常显示。调用setlocale(LC_ALL, “”)后会同时影响宽字符和窄字符，因此输出窄字符又呈现乱码，在后面调用setlocale(LC_ALL, “C”)还原默认的区域后窄字符能够正确显示。这段代码还不能在CDT的控制台中运行，只能在windows的控制台中测试，因为eclipse控制台采用的编码和windows控制台又不同，使用minGW的CDT无论编译还是输出对本地化的支持都不是很友好。相反在高版本的VS(Visual Studio)中很少碰到显示乱码问题，调用setlocale()设置宽字符区域不会影响窄字符的显示，将控制台设置为utf-8或ANSI都能正确显示，因为对于窄字符无论如何设置VS都会尝试将编码转化为控制台使用的编码输出。不过现在我们明白宽字符的运作原理了，那就是编译时将文本转为unicode-16或unicode-32的宽字符，运行时将unicode-16或unicode-32转化为ANSI或utf-8的窄字符。运行时是否能正确显示，是否有乱码要看编译器将宽字符或窄字符呈现为控制台显示编码的能力。宽字符是对传统字符的扩展，在处理上采用完全不同的方案，因此在使用上窄字符和宽字符相互独立，C代码中可以同时存在两套编码，互不影响，对应的处理函数也不同。
————————————————

