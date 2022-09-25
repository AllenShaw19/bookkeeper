package _default

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMmapFile(t *testing.T) {
	fmt.Println("test")
	f := &MmapFile{}
	assert.NotNil(t, f)
}
