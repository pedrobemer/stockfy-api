package order

import "stockfyApi/entity"

type Repository interface {
	Create(orderInsert entity.Order) entity.Order
}
