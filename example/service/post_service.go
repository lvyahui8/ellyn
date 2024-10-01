package service

import (
	"example/dao"
	"example/dao/model"
)

type PostService struct {
	postRepo dao.PostRepository
}

func (s *PostService) GetPostList(userId uint) (res []*model.Post) {
	posts := s.postRepo.FindByUser(userId)
	for _, post := range posts {
		res = append(res, &post)
	}
	return
}
