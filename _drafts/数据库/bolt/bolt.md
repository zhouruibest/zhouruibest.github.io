boltdb使用单内存映射文件作为存储（single memory-mapped file on disk）。boltdb在启动时会通过mmap系统调用将数据库文件映射到内存，这样可以仅通过内存访问来对文件进行读写，而将磁盘I/O交给操作系统管理，只有在事务提交或更新元数据时，boltdb才会通过fdatasyc系统调用强制将脏页落盘，以保证事务的ACID语义。

在linux系统中，内存与磁盘间的换入换出是以页为单位的。为了充分利用这一特定，boltdb的数据库文件也是按页组织的，且页大小与操作系统的页大小相等。

> 1. 为了在保证隔离性的同时支持“读读并发”、“读写并发”（boltdb不支持“写写并发”，即同一时刻只能有一个执行中的可写事务），boltdb在更新页时采用了Shadow Paging技术，其通过copy-on-write实现。在可写事务更新页时，boltdb首先会复制原页，然后在副本上更新，再将引用修改为新页上。这样，当可写事务更新页时，只读事务还可以读取原来的页；当创建读写事务时，boltdb会释放不再使用的页。这样，便实现了在支持“读读并发”、“读写并发”的同时保证事务的隔离性。

> 2. boltdb不会将空闲的页归还给系统。其原因有二：

1) 在不断增大的数据库中，被释放的页之后还会被重用。
2) boltdb为了保证读写并发的隔离性，使用copy-on-write来更新页，因此会在任意位置产生空闲页，而不只是在文件末尾产生空闲页

# Page，bolt中的页

page.go
```go
type page struct {
	id       pgid     // 页ID，单调递增
	flags    uint16   // 页标志， 用来表示页面的类型
	count    uint16   // 页面中元素的个数
	overflow uint32   // 溢出页个数，当单页无法容纳数据时，可以用与该页相邻的页面保存溢出的数据
	ptr      uintptr  // 页的数据的起始位置
}
```

flags: boltdb中的页共有三种用途：保存数据库的元数据（meta page）1、保存空闲页列表(freelist page)、保存数据，因为boltdb中数据是按照B+树组织的，因此保存数据的页又可分为分支节点（branch page）和叶子节点（leaf page）两种.

## meta page

```go
type meta struct {
	magic    uint32 // 魔数
	version  uint32 // 用来标识该文件采用的数据库版本号
	pageSize uint32 // 用来标识文件采用的页大小
	flags    uint32 // 保留字段
	root     bucket // 根bucket的结构体
	freelist pgid   // 空闲页列表的首页ID
	pgid     pgid   // 下一个分配的页ID，即当前最大页ID+1，用于mmap扩容时为新页编号
	txid     txid   // 下一个事务的ID，单调递增
	checksum uint64
}
```

## branch page & leaf page

branch page与leaf page是boltdb中用来保存B+树节点的页。B+树的分支节点仅用来保存索引（key），而叶子节点既保存索引，又保存值（value）。boltdb支持任意长度的key和value，因此无法直接结构化保存key和value的列表。为了解决这一问题，branch page和leaf page的Page Body起始处是一个由定长的索引（branchPageElement或leafPageElement）组成的列表，第$i$个索引记录了第$i$个key或key/value的起始位置与key的长度或key/value各自的长度：

![branch-page结构示意图.svg](./branch-page结构示意图.svg)

```go

// branchPageElement represents a node on a branch page.
type branchPageElement struct {
	pos   uint32
	ksize uint32
	pgid  pgid
}

// key returns a byte slice of the node key.
func (n *branchPageElement) key() []byte {
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	return (*[maxAllocSize]byte)(unsafe.Pointer(&buf[n.pos]))[:n.ksize]
}
```

![leaf-page结构示意图.svg](./leaf-page结构示意图.svg)

