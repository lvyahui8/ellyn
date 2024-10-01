package resource

import (
	"context"
	"example/dao/model"
	"example/service"
	"fmt"
)

type HomeResource struct {
	userService service.UserService
	postService service.PostService
}

func (r HomeResource) MyProfile(ctx context.Context) (*model.User, []*model.Post) {
	uid := ctx.Value("uid")
	if uid == nil {
		fmt.Printf("invalid param[uid]")
		return nil, nil
	}
	user := r.userService.GetUser(uid.(uint))
	postList := r.postService.GetPostList(uid.(uint))
	return user, postList
}
