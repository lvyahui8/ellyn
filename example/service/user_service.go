package service

import (
	"example/dao"
	"example/dao/model"
	"fmt"
)

type UserService struct {
	userRepo dao.UserRepository
}

func (s UserService) GetUser(id uint) *model.User {
	find, err := s.userRepo.Find(id)
	if err != nil {
		fmt.Printf("user not found. err %v", err)
		return nil
	}
	return &find
}
