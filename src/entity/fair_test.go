package entity_test

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/marcovargas74/m74-fair-api/src/entity"
)

func TestNewFair(t *testing.T) {
	fair, err := entity.NewFair("PRACA SANTA HELENA", "VILA PRUDENTE", "Leste", "VL ZELINA")
	assert.Equal(t, err, nil)
	assert.Equal(t, fair.Name, "PRACA SANTA HELENA")
}

func TestFairValidate(t *testing.T) {

	tests := []struct {
		give         string
		name         string
		district     string
		region5      string
		neighborhood string
		want         error
	}{
		{
			give:         "TestFairValidate Test1",
			name:         "minha feira",
			district:     "logo ali",
			region5:      "OESTE",
			neighborhood: "Meu Bairro",
			want:         nil,
		},
		{
			give:         "TestFairValidate Invalid Name",
			name:         "",
			district:     "logo ali",
			region5:      "OESTE",
			neighborhood: "Meu Bairro",
			want:         entity.ErrInvalidEntity,
		},
		{
			give:         "TestFairValidate Invalid District",
			name:         "outra feita",
			district:     "di",
			region5:      "OESTE",
			neighborhood: "Meu Bairro",
			want:         entity.ErrInvalidEntity,
		},
		{
			give:         "TestFairValidate Invalid region5",
			name:         "outra feita",
			district:     "logo ali",
			region5:      "NOROESTE",
			neighborhood: "Meu Bairro",
			want:         entity.ErrInvalidEntity,
		},
		{
			give:         "TestFairValidate Invalid neighborhood",
			name:         "outra feita",
			district:     "logo ali",
			region5:      "NOROESTE",
			neighborhood: "neighborhoodneighborhoodneighborhood",
			want:         entity.ErrInvalidEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			_, err := entity.NewFair(tt.name, tt.district, tt.region5, tt.neighborhood)
			assert.Equal(t, err, tt.want)
		})

	}

}
