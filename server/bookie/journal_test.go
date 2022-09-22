package bookie

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

// Test Read
func read0(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return string(content)
}

func read1(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func read2(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	chunks := make([]byte, 0)
	buf := make([]byte, 1024) //一次读取多少个字节
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf[:n]...)
	}
	return string(chunks)
}

func read3(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	chunks := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf[:n]...)
	}
	return string(chunks)
}

func TestRead(t *testing.T) {

	file := "./tmp/log.txt"

	start := time.Now()

	read0(file)
	t0 := time.Now()
	fmt.Printf("Cost time %v\n", t0.Sub(start))

	read1(file)
	t1 := time.Now()
	fmt.Printf("Cost time %v\n", t1.Sub(t0))

	read2(file)
	t2 := time.Now()
	fmt.Printf("Cost time %v\n", t2.Sub(t1))

	read3(file)
	t3 := time.Now()
	fmt.Printf("Cost time %v\n", t3.Sub(t2))

}

// Test Write
// ioutil.WriteFile
func write0(path string) {
	content := []byte("测试1\n测试2\n")
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		panic(err)
	}
}

// io.Write
func write1(path string) {
	content := "测试1\n测试2\n"

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModePerm) //打开文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err1 := io.WriteString(f, content) //写入文件(字符串)
	if err1 != nil {
		panic(err1)
	}
}

// file.Write
func write2(path string) {
	content := []byte("测试1\n测试2\n")

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm) //打开文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write(content) //写入文件(字节数组)
	if err != nil {
		panic(err)
	}
	f.Sync()
}

//
func write3(path string) {
	content := []byte("测试1\n测试2\n")

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModePerm) //打开文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f) //创建新的 Writer 对象
	_, err = w.Write(content)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func TestWrite(t *testing.T) {
	file := "./tmp/log.txt"
	start := time.Now()

	write0(file)
	t0 := time.Now()
	fmt.Printf("Cost time %v\n", t0.Sub(start))

	write1(file)
	t1 := time.Now()
	fmt.Printf("Cost time %v\n", t1.Sub(t0))

	write2(file)
	t2 := time.Now()
	fmt.Printf("Cost time %v\n", t2.Sub(t1))

	write3(file)
	t3 := time.Now()
	fmt.Printf("Cost time %v\n", t3.Sub(t2))
}

// BenchmarkWrite0-8   	   21103	     55208 ns/op
// BenchmarkWrite1-8   	   57613	     19406 ns/op
// BenchmarkWrite2-8   	      74	  17422775 ns/op
// BenchmarkWrite3-8   	   54492	     19836 ns/op
// BenchmarkMmap-8   	   30228	     49058 ns/op
func BenchmarkWrite0(b *testing.B) {
	file := "./tmp/log.txt"
	for i := 0; i < b.N; i++ {
		write0(file)
	}
}
func BenchmarkWrite1(b *testing.B) {
	file := "./tmp/log.txt"
	for i := 0; i < b.N; i++ {
		write1(file)
	}
}
func BenchmarkWrite2(b *testing.B) {
	file := "./tmp/log.txt"
	for i := 0; i < b.N; i++ {
		write2(file)
	}
}
func BenchmarkWrite3(b *testing.B) {
	file := "./tmp/log.txt"
	for i := 0; i < b.N; i++ {
		write3(file)
	}
}

const defaultMaxFileSize = 1 << 30        // 假设文件最大为 1G
const defaultMemMapSize = 128 * (1 << 20) // 假设映射的内存大小为 128M

type Demo struct {
	file    *os.File
	data    *[defaultMaxFileSize]byte
	dataRef []byte
}

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf(msg, v...))
	}
}

func (demo *Demo) mmap() {
	b, err := syscall.Mmap(int(demo.file.Fd()), 0, defaultMemMapSize, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	_assert(err == nil, "failed to mmap %v", err)
	demo.dataRef = b
	demo.data = (*[defaultMaxFileSize]byte)(unsafe.Pointer(&b[0]))
}

func (demo *Demo) grow(size int64) {
	if info, _ := demo.file.Stat(); info.Size() >= size {
		return
	}
	_assert(demo.file.Truncate(size) == nil, "failed to truncate")
}

func (demo *Demo) munmap() {
	_assert(syscall.Munmap(demo.dataRef) == nil, "failed to munmap")
	demo.data = nil
	demo.dataRef = nil
}

func mmapWrite(path string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm) //打开文件
	if err != nil {
		panic(err)
	}

	demo := &Demo{file: f}
	demo.grow(1)
	demo.mmap()
	defer demo.munmap()

	content := []byte("测试1\n测试2\n")

	demo.grow(int64(len(content) * 2))
	for i, v := range content {
		demo.data[i] = v
	}
}

// BenchmarkMmap-8   	   30228	     49058 ns/op
func BenchmarkMmap(b *testing.B) {
	file := "./tmp/log.txt"
	for i := 0; i < b.N; i++ {
		mmapWrite(file)
	}
}

// BenchmarkRead0-8   	  236522	      5052 ns/op
// BenchmarkRead1-8   	  218145	      5485 ns/op
// BenchmarkRead2-8   	  220753	      5487 ns/op
// BenchmarkRead3-8   	  238801	      4990 ns/op
func BenchmarkRead0(b *testing.B) {
	file := "./tmp/log.txt"
	for i := 0; i < b.N; i++ {
		read0(file)
	}
}
func BenchmarkRead1(b *testing.B) {
	file := "./tmp/log.txt"
	for i := 0; i < b.N; i++ {
		read1(file)
	}
}
func BenchmarkRead2(b *testing.B) {
	file := "./tmp/log.txt"
	for i := 0; i < b.N; i++ {
		read2(file)
	}
}
func BenchmarkRead3(b *testing.B) {
	file := "./tmp/log.txt"
	for i := 0; i < b.N; i++ {
		read3(file)
	}
}
