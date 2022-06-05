package service

import (
	"errors"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

var commentCtl = repository.GetCommentCtl()

func Publish(comment *model.Comment) error {
	return commentCtl.Publish(comment)
}

func DeleteById(userId, commentId int64) error {

	comment, err := commentCtl.QueryById(commentId)
	if err != nil {
		return err
	}

	if userId != comment.User.ID {
		return errors.New("用户无删除权限")
	}

	return commentCtl.DeleteById(commentId)
}

func QueryListByVideoId(videoId int64) ([]*model.Comment, error) {
	return commentCtl.QueryListByVideoId(videoId)
}
