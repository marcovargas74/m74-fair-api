package logs

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestThisFunction(t *testing.T) {
	assert.Equal(t, ThisFunction(), "logs.TestThisFunction(10)")
}

func TestChopPath(t *testing.T) {
	assert.Equal(t, chopPath("aDirName/fileName"), "fileName")
	assert.Equal(t, chopPath("aDirName/"), "")
	assert.Equal(t, chopPath("aDirName/aDirName/aDirName/aDirName/123"), "123")
	assert.Equal(t, chopPath("/123"), "123")
}
