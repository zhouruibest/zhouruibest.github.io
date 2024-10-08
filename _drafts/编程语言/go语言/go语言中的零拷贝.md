相信那些曾经使用 Go 写过 proxy server 的同学应该对io.Copy()/io.CopyN()/io.CopyBuffer()/io.ReaderFrom 等接口和方法不陌生，它们是使用 Go 操作各类 I/O 进行数据传输经常需要使用到的 API，其中基于 TCP 协议的 socket 在使用上述接口和方法进行数据传输时利用到了 Linux 的零拷贝技术 sendfile 和 splice。

https://strikefreedom.top/archives/pipe-pool-for-splice-in-go