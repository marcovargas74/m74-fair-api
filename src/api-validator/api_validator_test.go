package m74validatorapi

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestVersion(t *testing.T) {

	tests := []struct {
		give      string
		wantValue string
	}{
		{
			give:      "Test if get version OK",
			wantValue: "2022",
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			version := Version()
			assert.Equal(t, version[0:4], tt.wantValue)
		})

	}

}
