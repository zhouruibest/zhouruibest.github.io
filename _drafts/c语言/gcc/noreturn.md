`__attribute__ ((__noreturn__))` 是 GCC 中的函数属性，用于指示某个函数不会返回。具体来说，这个属性告诉编译器，在函数执行完毕后，程序将无法继续执行到该函数被调用的位置。通常情况下，这种情况都是由于函数内部抛出了异常、调用了 `exit()` 或 `abort()` 等函数，或者使用了死循环等造成的。

在 GCC 中，加上 `__attribute__ ((__noreturn__))` 属性的函数，编译器将自动优化代码，以避免对该函数进行不必要的处理。这可以帮助提高代码优化和提高程序性能。

```c
void fatal_error(char* message) __attribute__ ((__noreturn__));


#ifdef __GNUC__

    /* This macro stops 'gcc -Wall' complaining that "control reaches
       end of non-void function" if we use the following functions to
       terminate main() or some other non-void function. */

#define NORETURN __attribute__ ((__noreturn__))
#else
#define NORETURN
#endif
```

