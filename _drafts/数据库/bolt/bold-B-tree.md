# blot中的B+树
boltdb的B+Tree节点实现可分为存储中的实现（mmap memory）与内存中的实现（heap memory）两部分。

虽然boltdb通过mmap的方式将数据库文件映射到了内存中，但boltdb不会直接修改mmap的内存空间，而是只读mmap内存空间。当需要更新B+Tree的节点时，boltdb会读取mmap内存中相应的page，并在heap memory中构建相应的数据结构来修改，最后再通过pwrite+fdatasync的方式写入底层文件。

B+Tree节点内存部分主要由node结构体实现。boltdb中node是按需实例化的，对于不需要修改的node，boltdb直接从page中读取数据；而当boltdb需要修改B+Tree的某个节点时，则会将该节点从page实例化为node。在修改node时，boltdb会为其分配page buffer（dirty page），等到事务提交时，才会将这些page buffer中的数据统一落盘。

```go
// node represents an in-memory, deserialized page.
type node struct {
	bucket     *Bucket // 该node所属的Bucket节点
	isLeaf     bool    // 该node是否为叶子结点
	unbalanced bool    // 当前node是否可能不平衡
	spilled    bool    // 当前node是否已经被调整过
	key        []byte  // 保存node初始化时的第一个key
	pgid       pgid    // 当前node在mmap中的页id
	parent     *node   // 父节点指针
	children   nodes   // 保存已经实例化的孩子节点的指针，用于spill时递归向下更新node
	inodes     inodes  // 该node的内部节点，即该node包含的元素
}
```