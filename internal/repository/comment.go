package repository

import (
	"time"

	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

// comment struct mapped to database
type Comment struct {
	ID          uint      `gorm:"primarykey"`
	userId      int64     `gorm:"user_id"`
	videoId     int64     `gorm:"video_id"`
	commentText string    `gorm:"comment_text"`
	CreatedAt   time.Time `gorm:"create_date"`
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type CommentCtl struct{}

var commentCtl CommentCtl

// 指定 Comment 对应的数据库表名
// func (Comment) TableName() string {
// 	return "comments"
// }

func GetCommentCtl() model.CommentCtl {
	return &commentCtl
}

func (commentCtl *CommentCtl) Publish(comment *model.Comment) error {
	return dbProvider.GetDB().Create(comment).Error
}

func (commentCtl *CommentCtl) DeleteById(commentId int64) error {
	return dbProvider.GetDB().Where("id = ?", commentId).Delete(&Comment{}).Error
}

func (commentCtl *CommentCtl) QueryListByVideoId(videoId int64) ([]model.Comment, error) {

	var commentEntitis []Comment
	var comments []model.Comment

	err := dbProvider.GetDB().Where("video_id = ?", videoId).Find(&commentEntitis).Error

	return comments, err
}
