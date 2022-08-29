package logs

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestThisFunction(t *testing.T) {
	//fair, err := entity.NewFair("PRACA SANTA HELENA", "VILA PRUDENTE", "Leste", "VL ZELINA")
	//assert.Equal(t, err, nil)
	assert.Equal(t, ThisFunction(), "logs.TestThisFunction(12)")
}

func TestChopPath(t *testing.T) {
	//fair, err := entity.NewFair("PRACA SANTA HELENA", "VILA PRUDENTE", "Leste", "VL ZELINA")
	//assert.Equal(t, err, nil)
	assert.Equal(t, chopPath("aDirName/fileName"), "fileName")
	assert.Equal(t, chopPath("aDirName/"), "")
	assert.Equal(t, chopPath("aDirName/aDirName/aDirName/aDirName/123"), "123")
	assert.Equal(t, chopPath("/123"), "123")
}
