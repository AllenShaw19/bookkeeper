package _default

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)


var (
	fileName = "test.file"
)
func TestNewMmapFile(t *testing.T) {
	fmt.Println("test")
	f, err := NewMmapFile(fileName, defaultMemMapSize)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.NotNil(t, f.wrotePos)

	n := 10
	for i := 0; i < n; i++ {
		data := []byte(fmt.Sprintf("test test data %d\n", i))
		f.AppendMessage(data)
	}

	err = f.Close()

}
