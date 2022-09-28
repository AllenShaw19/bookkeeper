package _default

import (
	"go.uber.org/atomic"
	"os"
	"syscall"
	"unsafe"
)

const (
	OsPageSize = 1024 * 4
	defaultMaxFileSize = 1 << 30        // 假设文件最大为 1G
	defaultMemMapSize = 128 * (1 << 20) // 假设映射的内存大小为 128M
)

var (
	TotalMappedVirtualMemory atomic.Int64
	TotalMappedFile          atomic.Int32
)

type MmapFile struct {
	startPos       atomic.Int32
	wrotePos       *atomic.Int32
	committedPos   atomic.Int32
	flushedPos     atomic.Int32
	fileSize       int
	fileFromOffset int64
	fileName       string
	storeTimestamp atomic.Int64

	// mmapfile
	file           *os.File
	data *[defaultMaxFileSize]byte
	offset int
	dataRef []byte
}

func NewMmapFile(fileName string, fileSize int) (*MmapFile, error) {
	mf := &MmapFile{}
	mf.fileName = fileName
	mf.fileSize = fileSize
	mf.wrotePos = atomic.NewInt32(0)
	ok := false

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		//
		return nil, err
	}
	defer func() {
		if !ok && mf.file != nil {
			mf.file.Close()
		}
	}()

	mf.file = f
	buff, err := syscall.Mmap(int(f.Fd()), 0, fileSize, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		//
		return nil, err
	}

	mf.dataRef = buff
	mf.data = (*[defaultMaxFileSize]byte)(unsafe.Pointer(&buff[0]))

	TotalMappedVirtualMemory.Add(int64(fileSize))
	TotalMappedFile.Inc()
	ok = true

	return mf, nil
}

func (f *MmapFile) grow(size int64) error {
	stat, err := f.file.Stat()
	if err != nil {
		//
		return err
	}
	if stat.Size() >= size {
		return nil
	}
	err = f.file.Truncate(size)
	if err != nil {
		//
		return err
	}
	return nil
}

func (f *MmapFile) AppendMessage(data []byte) bool {
	currPos := int(f.wrotePos.Load())

	length := len(data)
	if currPos+length <= f.fileSize {
		err := f.grow(int64(currPos + length))
		if err != nil {
			return false
		}
		copy(f.data[currPos:], data)
		f.wrotePos.Add(int32(length))
		return true
	}
	return false
}

func (f *MmapFile) Close() error {
	err := syscall.Munmap(f.dataRef)
	if err != nil {
		return err
	}

	f.data = nil
	f.dataRef = nil
	return nil
}