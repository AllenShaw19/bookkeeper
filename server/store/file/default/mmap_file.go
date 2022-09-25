package _default

import (
	"bytes"
	"go.uber.org/atomic"
	"os"
	"syscall"
)

const (
	OsPageSize = 1024 * 4
)

var (
	TotalMappedVirtualMemory atomic.Int64
	TotalMappedFile          atomic.Int32
)

type MmapFile struct {
	startPos       atomic.Int32
	wrotePos       atomic.Int32
	committedPos   atomic.Int32
	flushedPos     atomic.Int32
	file           *os.File
	fileSize       int
	fileFromOffset int64
	fileName       string
	mappedByteBuff *bytes.Buffer
	storeTimestamp atomic.Int64
}

func NewMmapFile(fileName string, fileSize int) (*MmapFile, error) {
	mf := &MmapFile{}
	mf.fileName = fileName
	mf.fileSize = fileSize
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
	buff, err := syscall.Mmap(int(f.Fd()), 0, fileSize, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_FILE)
	if err != nil {
		//
		return nil, err
	}
	mf.mappedByteBuff = bytes.NewBuffer(buff)

	TotalMappedVirtualMemory.Add(int64(fileSize))
	TotalMappedFile.Inc()
	ok = true

	return mf, nil
}

func (f *MmapFile) AppendMessage(data []byte, offset, length int) bool {
	currPos := int(f.wrotePos.Load())

	if currPos+length <= f.fileSize {
	}
}
