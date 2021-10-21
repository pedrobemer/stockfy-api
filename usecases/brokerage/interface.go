package brokerage

import "stockfyApi/entity"

type Repository interface {
	Search(specificFetch string, args ...string) ([]entity.Brokerage, error)
}

type UseCases interface {
	SearchBrokerage(searchType string, name string, country string) (
		[]entity.Brokerage, error)
}
