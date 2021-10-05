package brokerage

import "stockfyApi/entity"

type Repository interface {
	Search(specificFetch string, args ...string) ([]entity.Brokerage, error)
}
