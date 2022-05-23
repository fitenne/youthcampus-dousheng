package service

import (
	"errors"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"log"
	"strconv"
)

var commentCtl = repository.GetCommentCtl()

func Publish(videoId int64, comment *model.Comment) error {
	return commentCtl.Publish(videoId, comment)
}

func DeleteById(userId, commentId int64) error {
	comment, err := commentCtl.QueryById(commentId)
	if err != nil {
		log.Println(err.Error())
		return errors.Unwrap(err)
	}

	if comment.User.ID == 0 {
		return errors.New("指定id不存在记录")
	}

	if userId != comment.User.ID {
		return errors.New("该用户没有权限：" + strconv.Itoa(int(userId)))
	}

	return commentCtl.DeleteById(commentId)
}

func QueryListByVideoId(videoId int64) ([]model.Comment, error) {
	return commentCtl.QueryListByVideoId(videoId)
}
