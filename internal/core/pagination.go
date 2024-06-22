package core

const itemsPerPage = 50

func ItemsPerPage() int {
	return itemsPerPage
}

type PaginationParams struct {
	Page int
}
