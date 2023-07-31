```c
#include <stdio.h>
#include <string.h>

#define STR "THISIS1234"

int main(void) {
    char s[40] = "helloworld";
    char *s2 = "hehehe";
    char s3[] = "hehehhe";
    printf("s    sizeof = %d, strlen = %d \n", sizeof s, strlen(s));
    printf("*s   sizeof = %d \n", sizeof *s);
    printf("s2   sizeof = %d, strlen = %d \n", sizeof s2, strlen(s2));
    printf("*s2  sizeof = %d \n", sizeof *s2);
    printf("s3   sizeof = %d, strlen = %d \n", sizeof s3, strlen(s3));
    printf("*s3  sizeof = %d \n", sizeof *s3);
    printf("STR  sizeof = %d, strlen = %d \n", sizeof STR, strlen(STR));
    return 0;
}
```

输出

```sh
s    sizeof = 40, strlen = 10
*s   sizeof = 1
s2   sizeof = 8, strlen = 6
*s2  sizeof = 1
s3   sizeof = 8, strlen = 7
*s3  sizeof = 1
STR  sizeof = 11, strlen = 10
```