package dao

import (
	"example/dao/model"
	"fmt"
)

type UserRepository struct {
}

func (u UserRepository) Find(id uint) (model.User, error) {
	return model.User{
		Name: fmt.Sprintf("uName%d", id),
		Age:  int(id % 20),
	}, nil
}

func (u UserRepository) Insert(t model.User) error {
	return nil
}
