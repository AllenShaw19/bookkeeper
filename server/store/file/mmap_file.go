package file

import "os"

type File interface {
	GetFileName() string
	GetFileSize() int64
	GetFile() *os.File
	IsFull() bool
	IsAvailable() bool
	AppendMessage(data []byte) error
	AppendMessageWithOffset(data []byte, offset, length int) error
	GetFileFromOffset() int64
	Flush() int64
	Commit() int64
}
