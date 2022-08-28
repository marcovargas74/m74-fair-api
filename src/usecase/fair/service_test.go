package fair

import (
	"fmt"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/marcovargas74/m74-fair-api/src/entity"
)

func newFaikeFair() *entity.Fair {
	return &entity.Fair{
		Name:         "PRACA SANTA HELENA",
		District:     "VILA PRUDENTE",
		Region5:      "Leste",
		Neighborhood: "VL ZELINA",
		CreatedAt:    time.Now(),
	}
}

func Test_Create(t *testing.T) {

	repo := NewInmem()
	server := NewService(repo)
	aFair := newFaikeFair()
	_, err := server.CreateFair(aFair.Name, aFair.District, aFair.Region5, aFair.Neighborhood)

	assert.Equal(t, err, nil)
	assert.Equal(t, aFair.CreatedAt.IsZero(), false)
}

func Test_SearchAndFind(t *testing.T) {
	repo := NewInmem()
	server := NewService(repo)
	aFair1 := newFaikeFair()
	aFair2 := newFaikeFair()
	aFair2.Name = "Feira da Madrugada"

	uID, _ := server.CreateFair(aFair1.Name, aFair1.District, aFair1.Region5, aFair1.Neighborhood)
	_, _ = server.CreateFair(aFair2.Name, aFair2.District, aFair2.Region5, aFair2.Neighborhood)

	t.Run("search", func(t *testing.T) {
		result, err := server.SearchFairs("name", "praca")
		assert.Equal(t, err, nil)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "PRACA SANTA HELENA", result[0].Name)

		_, err = server.SearchFairs("region5", "nosul")
		assert.Equal(t, entity.ErrNotFound, err)
	})

	t.Run("list all", func(t *testing.T) {
		all, err := server.ListFairs()
		assert.Equal(t, err, nil)
		assert.Equal(t, 2, len(all))
	})

	t.Run("get", func(t *testing.T) {
		saved, err := server.GetFair(uID)
		assert.Equal(t, err, nil)
		assert.Equal(t, aFair1.Name, saved.Name)
	})
}

func Test_Update(t *testing.T) {
	repo := NewInmem()
	server := NewService(repo)
	aFair := newFaikeFair()

	id, err := server.CreateFair(aFair.Name, aFair.District, aFair.Region5, aFair.Neighborhood)
	assert.Equal(t, err, nil)
	saved, _ := server.GetFair(id)

	saved.Name = "Praça"
	err = server.UpdateFair(saved)
	fmt.Printf("UPDATE erro[%s]\n", err)
	updated, err := server.GetFair(id)
	fmt.Printf("UPDATE erro[%s]saved[%v]\n", err, updated)
	assert.Equal(t, "Praça", updated.Name)
}

func TestDelete(t *testing.T) {
	repo := NewInmem()
	server := NewService(repo)
	aFair1 := newFaikeFair()
	aFair2 := newFaikeFair()
	aFair2ID, _ := server.CreateFair(aFair2.Name, aFair2.District, aFair2.Region5, aFair2.Neighborhood)

	err := server.DeleteFair(aFair1.ID)
	assert.Equal(t, entity.ErrNotFound, err)

	err = server.DeleteFair(aFair2ID)
	assert.Equal(t, err, nil)
	_, err = server.GetFair(aFair2ID)
	assert.Equal(t, entity.ErrNotFound, err)
}
