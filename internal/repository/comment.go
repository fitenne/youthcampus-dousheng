package repository

import (
	"errors"
	"gorm.io/gorm"

	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

type CommentCtl struct{}

var commentCtl CommentCtl

func GetCommentCtl() model.CommentCtl {
	return &commentCtl
}

// Publish 保存评论
func (commentCtl *CommentCtl) Publish(comment *model.Comment) error {

	return dbProvider.GetDB().Transaction(func(tx *gorm.DB) error {

		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Create(&comment).Error; err != nil {
			return err // 返回任何错误都会回滚事务
		}

		if err := tx.Table("videos").Where("id", comment.VideoId).
			UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).
			Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

// DeleteById 删除评论(软删除)
func (commentCtl *CommentCtl) DeleteById(commentId, videoId int64) error {

	return dbProvider.GetDB().Transaction(func(tx *gorm.DB) error {

		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Where("id = ?", commentId).Delete(&model.Comment{}).Error; err != nil {
			return err // 返回任何错误都会回滚事务
		}

		if err := tx.Table("videos").Where("id", videoId).
			UpdateColumn("comment_count", gorm.Expr("comment_count - 1")).
			Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
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
		Order("created_at desc").
		Find(&comments).Error

	return comments, err
}
