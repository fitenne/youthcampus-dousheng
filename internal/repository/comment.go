package repository

import (
	"errors"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

type CommentCtl struct{}

var commentCtl CommentCtl

func GetCommentCtl() model.CommentCtl {
	return &commentCtl
}

// Publish 保存评论
func (commentCtl *CommentCtl) Publish(comment *model.Comment) error {
	return dbProvider.GetDB().Create(&comment).Error
}

// DeleteById 删除评论(软删除)
func (commentCtl *CommentCtl) DeleteById(commentId int64) error {
	return dbProvider.GetDB().Where("id = ?", commentId).Delete(&model.Comment{}).Error
}

// QueryById 查询评论
func (commentCtl *CommentCtl) QueryById(commentId int64) (*model.Comment, error) {

	// 单个查询
	comment := model.Comment{}

	// 异常处理
	err := dbProvider.GetDB().Preload("User").Find(&comment, commentId).Error
	if err != nil {
		return nil, err
	}

	if comment.ID <= 0 {
		return nil, errors.New("comment 不存在")
	}

	// 返回结果
	return &comment, err
}

// QueryListByVideoId 列表查询
func (commentCtl *CommentCtl) QueryListByVideoId(videoId int64) ([]model.Comment, error) {

	// 数据库实体
	comments := make([]model.Comment, 20)

	// 查库
	err := dbProvider.GetDB().
		Where("video_id = ?", videoId).
		Joins("User").
		Order("create_date desc").
		Find(&comments).Error

	return comments, err
}
