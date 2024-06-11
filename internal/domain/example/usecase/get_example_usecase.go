package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/example/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *ExampleService) GetExample(id string) (*entity.Example, error) {
	example, err := s.exampleRepo.FindByID(id)
	if err != nil {
		return nil, validation.NewNotFoundErr("example")
	}

	return example, nil
}
