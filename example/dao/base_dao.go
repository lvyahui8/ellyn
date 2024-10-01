package dao

import "example/dao/model"

type Repository[T model.Model] interface {
	Find(id uint) (T, error)
	Insert(t T) error
}
