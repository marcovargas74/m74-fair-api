package fair

import (
	"strings"
	"time"

	"github.com/marcovargas74/m74-fair-api/src/entity"
	"github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

//Service fair usecase
type Service struct {
	repo Repository
}

//NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

//CreateFair create a Fair
func (s *Service) CreateFair(name string, district string, region5 string, neighborhood string) (entity.ID, error) {
	f, err := entity.NewFair(name, district, region5, neighborhood)
	if err != nil {
		return entity.NewID(), err
	}
	return s.repo.Create(f)
}

//GetFair get a Fair
func (s *Service) GetFair(id entity.ID) (*entity.Fair, error) {

	logs.Debug("getFair(1.SERVICE)ID %s", id)
	b, err := s.repo.Get(id)
	if b == nil {
		logs.Error("getFairERROR(1.SERVICE)ID %s", id)
		return nil, entity.ErrNotFound
	}
	if err != nil {
		logs.Error("getFairERROR(1.SERVICE)ID %s", id)
		return nil, err
	}

	logs.Debug("getFair(OK.SERVICE)ID %s", id)
	return b, nil
}

//SearchFairs search Fair
func (s *Service) SearchFairs(key string, value string) ([]*entity.Fair, error) {
	fairs, err := s.repo.Search(strings.ToLower(key), strings.ToLower(value))
	if err != nil {
		return nil, err
	}
	if len(fairs) == 0 {
		return nil, entity.ErrNotFound
	}
	return fairs, nil
}

//ListFairs list Fairs
func (s *Service) ListFairs() ([]*entity.Fair, error) {
	fairs, err := s.repo.List()
	if err != nil {
		logs.Error("SEESSSSSER 11111 FALHA ao conectar ao Banco Mysql do DOcker %v", err)
		return nil, err
	}
	/*if len(fairs) == 0 {
		logs.Error("SEESSSSSER tabale vazia %v", err)
		return nil, entity.ErrNotFound
	}*/
	return fairs, nil
}

//DeleteFair Delete a Fair
func (s *Service) DeleteFair(id entity.ID) error {
	_, err := s.GetFair(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

//UpdateFair Update a Fair
func (s *Service) UpdateFair(e *entity.Fair) error {
	err := e.Validate()
	if err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return s.repo.Update(e)
}
