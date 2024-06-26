package core

import (
	"errors"
)

const itemsPerPage = 50

func ItemsPerPage() int {
	return itemsPerPage
}

func NewPaginationParams(page int) (*PaginationParams, error) {
	if page < 0 {
		return nil, errors.New("page must be a positive value")
	}

	return &PaginationParams{
		Page: page,
	}, nil
}

type PaginationParams struct {
	Page int
}
