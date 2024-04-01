package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const defaultMaxFileSize = 1 << 30        // 假设文件最大为 1G
const defaultMemMapSize = 128 * (1 << 20) // 假设映射的内存大小为 128M

func main() {
	mmpFile := NewMmpFile("test.txt")
	defer mmpFile.munmap()
	defer mmpFile.file.Close()
	msg := "hello csdn colinrs!"

	mmpFile.grow(int64(len(msg) * 2))
	for i, v := range msg {
		mmpFile.data[i] = byte(v)
	}
}

type MmpFile struct {
	file    *os.File
	data    *[defaultMaxFileSize]byte
	dataRef []byte
}

func NewMmpFile(fileName string) *MmpFile {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Open file error: %v", err)
	}
	mmpFile := &MmpFile{file: file}
	mmpFile.mmap()
	return mmpFile
}

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf(msg, v...))
	}
}

/*
- fd：待映射的文件描述符。
- offset：映射到内存区域的起始位置，0 表示由内核指定内存地址。
- length：要映射的内存区域的大小。
- prot：内存保护标志位，可以通过或运算符`|`组合
    - PROT_EXEC  // 页内容可以被执行
    - PROT_READ  // 页内容可以被读取
    - PROT_WRITE // 页可以被写入
    - PROT_NONE  // 页不可访问
- flags：映射对象的类型，常用的是以下两类
    - MAP_SHARED  // 共享映射，写入数据会复制回文件, 与映射该文件的其他进程共享。
    - MAP_PRIVATE // 建立一个写入时拷贝的私有映射，写入数据不影响原文件
**/

func (mmpFile *MmpFile) mmap() {
	b, err := syscall.Mmap(int(mmpFile.file.Fd()), 0,
		defaultMemMapSize, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	_assert(err == nil, "failed to mmap", err)
	mmpFile.dataRef = b
	mmpFile.data = (*[defaultMaxFileSize]byte)(unsafe.Pointer(&b[0]))
}

func (mmpFile *MmpFile) grow(size int64) {
	if info, _ := mmpFile.file.Stat(); info.Size() >= size {
		return
	}
	_assert(mmpFile.file.Truncate(size) == nil, "failed to truncate")
}

func (mmpFile *MmpFile) munmap() {
	_assert(syscall.Munmap(mmpFile.dataRef) == nil, "failed to munmap")
	mmpFile.data = nil
	mmpFile.dataRef = nil
}
