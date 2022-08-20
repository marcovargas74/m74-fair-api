package fair

import (
	"github.com/marcovargas74/m74-val-cpf-cnpj/src/entity"
)

//Reader interface
type Reader interface {
	Get(id entity.ID) (*entity.Fair, error)
	Search(query string) ([]*entity.Fair, error)
	List() ([]*entity.Fair, error)
}

//Writer book writer
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
	SearchFair(query string) ([]*entity.Fair, error)
	ListFairs() ([]*entity.Fair, error)
	CreateFair(name string, district string, region5 string, neighborhood string) (entity.ID, error)
	UpdateFair(e *entity.Fair) error
	DeleteFair(id entity.ID) error
}
