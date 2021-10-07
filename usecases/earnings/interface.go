package earnings

import "stockfyApi/entity"

type Repository interface {
	Create(earningOrder entity.Earnings) ([]entity.Earnings, error)
}
