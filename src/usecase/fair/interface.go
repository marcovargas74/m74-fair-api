package fair

import (
	"github.com/marcovargas74/m74-fair-api/src/entity"
)

//Reader Fair Read interface
type Reader interface {
	Get(id entity.ID) (*entity.Fair, error)
	Search(key string, value string) ([]*entity.Fair, error)
	List() ([]*entity.Fair, error)
}

//Writer Fair writer interface
type Writer interface {
	Create(e *entity.Fair) (entity.ID, error)
	Update(e *entity.Fair) error
	Delete(id entity.ID) error
}

//Repository interface
type Repository interface {
	Reader
	Writer
}

//UseCase interface
type UseCase interface {
	GetFair(id entity.ID) (*entity.Fair, error)
	SearchFairs(key string, value string) ([]*entity.Fair, error)
	ListFairs() ([]*entity.Fair, error)
	CreateFair(name string, district string, region5 string, neighborhood string) (entity.ID, error)
	UpdateFair(e *entity.Fair) error
	DeleteFair(id entity.ID) error
}
