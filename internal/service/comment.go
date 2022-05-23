package service

import (
	"errors"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"strconv"
)

var commentCtl = repository.GetCommentCtl()

func Publish(videoId int64, comment *model.Comment) error {
	return commentCtl.Publish(videoId, comment)
}

func DeleteById(userId, commentId int64) error {
	comment, err := commentCtl.QueryById(commentId)
	if err != nil {
		return errors.Unwrap(err)
	}

	if userId != comment.User.ID {
		return errors.New("该用户没有权限：" + strconv.Itoa(int(userId)))
	}

	return commentCtl.DeleteById(commentId)
}

func QueryListByVideoId(videoId int64) ([]model.Comment, error) {
	return commentCtl.QueryListByVideoId(videoId)
}
