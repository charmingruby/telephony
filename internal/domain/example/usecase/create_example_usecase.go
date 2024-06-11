package usecase

import (
	"github.com/charmingruby/telephony/internal/domain/example/dto"
	"github.com/charmingruby/telephony/internal/domain/example/entity"
	"github.com/charmingruby/telephony/internal/validation"
)

func (s *ExampleService) CreateExample(dto dto.CreateExampleDTO) error {
	example, err := entity.NewExample(dto.Name)
	if err != nil {
		return err
	}

	if err := s.exampleRepo.Store(example); err != nil {
		return validation.NewInternalErr()
	}

	return nil
}
