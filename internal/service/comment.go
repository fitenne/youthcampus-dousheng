package service

import (
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

var commentCtl = repository.GetCommentCtl()

func Publish(comment *model.Comment) error {
	return commentCtl.Publish(comment)
}

func DeleteById(commentId int64) error {
	return commentCtl.DeleteById(commentId)
}

func QueryListByVideoId(videoId int64) ([]model.Comment, error) {
	return QueryListByVideoId(videoId)
}
