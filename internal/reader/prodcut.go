package reader

import "github.com/naumov-andrey/linxdatacenter-assignment/internal/model"

type ProductReader interface {
	Read() (model.Product, error)
}
