package dao

import (
	"example/dao/model"
	"fmt"
)

type PostRepository struct {
}

func (p PostRepository) Find(id uint) (model.Post, error) {
	return model.Post{
		Title:   "test",
		Desc:    "description",
		Content: "test content",
	}, nil
}

func (p PostRepository) Insert(t model.Post) error {
	return nil
}

func (p PostRepository) FindByUser(uid uint) (res []model.Post) {
	for i := 0; i < 10; i++ {
		res = append(res, model.Post{
			Title:   fmt.Sprintf("t%d,uid %d", i, uid),
			Desc:    fmt.Sprintf("descrip %d", i),
			Content: fmt.Sprintf("content %d", i),
		})
	}
	return
}
