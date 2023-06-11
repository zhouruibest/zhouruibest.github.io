# TCP连接的关闭

在大多数情况下，TCP连接都是先关闭一个方向，此时另外一个方向还可以正常进行数据传输。当客户端发起连接中断，此时客户端不再往服务器写入数据，此时可能服务器端正在对客户端的最后报文进行处理，当完成这些处理后，服务器端把结果通过套接字写给客户端，我们说这个套接字的状态此时是“半关闭”的。

## 关闭连接的方式

### close函数（关闭两个方向）

```v
int close(int sockfd)
```
这个函数对已经连接的socket执行close操作，成功返回0，失败返回-1。

close函数会对socket的引用计数-1，一旦socket的引用计数被减为0，就会对socket进行彻底释放，并且会关闭TCP两个方向的数据流。

socket引用计数：由于socket可以被多个进程共享，比如通过fork产生子进程，那么socket的引用计数就会+1，调用一次close函数，socket引用计数就-1。

为了关闭两个方向的数据流，在数据接收方向，系统内核会将socket设置为不可读，任何读操作都会返回异常；在数据发送方向，系统内核尝试将发送缓冲区的数据发送给对端，并最后向对端发送一个FIN报文，接下来如果再对socket进行写操作会返回异常。

如果对端没有检测到socket已关闭，仍然继续发送报文，则会收到一个RST报文。如果向这个已经收到RST的socket执行写操作，内核会发出一个SIGPIPE信号给进程，该信号的默认行为是终止进程。

> 在UNIX/LINUX下，非阻塞模式SOCKET可以采用recv+MSG_PEEK的方式进行判断连接是否断开，其中MSG_PEEK保证了仅仅进行状态判断，而不影响数据接收
>
> 对于主动关闭的SOCKET, recv返回-1，而且errno被置为9（#define EBADF   9 /* Bad file number */）
> 或104 （#define ECONNRESET 104 /* Connection reset by peer */）
>
> 对于被动关闭的SOCKET,recv返回0，而且errno被置为11（#define EWOULDBLOCK EAGAIN /* Operation would block */）
>
> 对正常的SOCKET, 如果有接收数据，则返回>0, 否则返回-1，而且errno被置为11（#define EWOULDBLOCK EAGAIN /* Operation would block */）
> 因此对于简单的状态判断（不过多考虑异常情况），
>
>    recv返回>0，   正常
>    返回-1，而且errno被置为11  正常
>    其它情况    关闭
> -----------------------------------
> 如何在C语言中判断socket是否已经断开
> https://blog.51cto.com/u_15127557/4198614

###  shutdown函数

```c
int shutdown(int sockfd, int howto)
```

对已连接的socket执行shutdown操作，成功返回0，失败返回-1。

howto为设置选项，主要有3个：

1. SHUT_RD(0)：关闭连接的“读”这个方向，对该socket进行读操作直接返回 EOF，从数据角度来看，套接字上接收缓冲区已有的数据将被丢弃，如果再有新的数据流到达，会对数据进行 ACK，然后丢弃。也就是说，对端还是会接收到 ACK，但是在这种情况下根本不知道数据已经被丢弃了。
2. SHUT_WR(1)：关闭连接的“写”这个方向，在这种情况下，连接处于“半关闭”状态，此时，不管socket引用计数的值是多少，都会直接关闭连接的写方向。套接字上发送缓冲区已有的数据将被立即发送出去，并发送一个 FIN 报文给对端。应用程序如果对该套接字进行写操作会报错。（此时对端仍然可以发送数据，调用shutdown的一端也可以对收到的数据发送ACK）（如果有进程共享此socket，那么也会受到影响）
3. SHUT_RDWR(2)：相当于 SHUT_RD 和 SHUT_WR 操作各一次，关闭套接字的读和写两个方向。

## close和shutdown的区别

1. close 会关闭连接，并释放所有连接对应的资源，而 shutdown 并不会释放掉套接字和所有的资源。确切地说，close用来关闭套接字，将套接字描述符（或句柄）从内存清除，之后再也不能使用该套接字。应用程序关闭套接字后，与该套接字相关的连接和缓存也失去了意义，TCP协议会自动触发关闭连接的操作。 shutdown() 用来关闭连接，而不是套接字，不管调用多少次 shutdown()，套接字依然存在，直到调用close将套接字从内存清除。(即调用shutdown后，仍然需要调用close关闭socket). 调用close关闭套接字，或调用shutdown关闭输出流时，都会向对方发送FIN包，FIN 包表示数据传输完毕，计算机收到 FIN 包就知道不会再有数据传送过来了。

2. close存在引用计数的概念，并不一定导致该套接字不可用，而shutdown则不会管引用计数接使得该套接字不可用，如果有别的进程企图使用该套接字，将会受到影响。

3. close 的引用计数导致不一定会发出 FIN 结束报文，而 shutdown 则总是会发出 FIN 结束报文，这在我们打算关闭连接通知对端的时候，是非常重要的

![close和shutdown的区别](./close%E5%92%8Cshutdown%E7%9A%84%E5%8C%BA%E5%88%AB.awebp)

## 为什么直接调用exit(0)就可以完成FIN报文的发送？为什么不需要调用close或者shutdown呢？

在调用exit(0)后进程会退出，与进程相关的所有资源，文件，内存，信号等内核分配的资源都会被释放。在linux中，一切皆文件，本身socket就是一种文件类型，内核会为每一个打开的文件创建file结构并维护指向该结构的引用计数，每一个进程结构中都会维护本进程打开的文件数组，数组下标就是fd，内容就指向上面的file结构，而close做的事就是删除本进程打开的文件数组中指定的fd项，并把指向的file结构中的引用计数减一，等引用计数为0的时候，就会调用内部包含的文件操作close。

## 调用close后发生了什么？

场景：服务端通过close()主动关闭一个TCP连接，客户端通过read()获得了0（read返回值为0表示EOF，即对端发送了FIN包），调用close()关闭这个连接。
在TCP层面：服务器调用close()后，向客户端发送FIN，客户端回应FIN-ACK。服务器进入FIN-WAIT-2状态，客户端进入CLOSE-WAIT状态。 客户端调用close()后，向服务端发送FIN，服务端会用FIN-ACK。服务端进入TIME-WAIT状态，客户端直接进入CLOSE状态，连接结束。
如果客户端在获得read()==0read()==0read()==0后，仍然向服务端写入数据，则会收到一个RST报文，如果向这个已经收到RST的socket继续执行写操作，内核会发出一个SIGPIPE信号给进程，该信号的默认行为是终止进程。
如果客户端在获得read()==0read()==0read()==0后，没有及时地调用close()，比如当read()==0read()==0read()==0时，客户端阻塞，等待一段时间后，如果发送SIGINT使进程退出（相当于调用close()，这种情况相当于没有及时调用close），会发生以下现象:

- 服务端SOCKET处于FIN-WAIT-2状态时，发送SIGINT信号使客户端退出，客户端发送FIN，服务端回复FIN-ACK。此时，按正常流程结束链接。
- 服务端SOCKET等待FIN-WAIT-2状态超时后，客户端发送FIN，服务端回复RST结束连接。

分析：服务端SOCKET关闭后，没有对这一SOCKET的引用。这一SOCKET进入到“孤儿SOCKET“的状态。孤儿Socket存在时，系统协议栈负责完成后续的FIN流程，当孤儿Socket超时后，系统协议栈将不存在这一Socket的信息。客户端此时发送FIN，将收到RST应答。

孤儿SOCKET：从应用程序来看，此条socket连接已经收发数据完毕，关闭了此连接，但是linux内核中为了完成正常的tcp协议（比如缓冲区中的数据）转换，会在内核的tcp协议层继续维护这些socket状态，直至系统回收。处于此种状态下的socket就是orphan socket（孤儿socket

