package repository

import (
	"gorm.io/gorm"
	"time"

	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

// comment struct mapped to database
type Comment struct {
	ID int64 `gorm:"primarykey"`
	// 使用 UserId 作为外键
	UserId    int64  `gorm:"user_id"`
	User      User   `gorm:"foreignKey:UserId"`
	VideoId   int64  `gorm:"video_id"`
	Content   string `gorm:"content"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`
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

func (commentCtl *CommentCtl) Publish(videoId int64, comment *model.Comment) error {

	// 创建评论信息
	err := dbProvider.GetDB().Create(&Comment{
		UserId:  comment.User.ID,
		VideoId: videoId,
		Content: comment.CommentText,
	}).Error

	return err
}

func (commentCtl *CommentCtl) DeleteById(commentId int64) error {
	// 软删除
	return dbProvider.GetDB().Where("id = ?", commentId).Delete(&Comment{}).Error
}

func (commentCtl *CommentCtl) QueryById(commentId int64) (*model.Comment, error) {
	// 单个查询
	var commentEntity Comment = Comment{}

	// 异常处理
	err := dbProvider.GetDB().Preload("User").Find(&commentEntity, commentId).Error
	if err != nil {
		return nil, err
	}

	// 返回结果
	return &model.Comment{
		ID: commentEntity.ID,
		User: model.User{
			ID:            commentEntity.User.ID,
			Name:          commentEntity.User.UserName,
			FollowCount:   commentEntity.User.FollowCount,
			FollowerCount: commentEntity.User.FollowerCount,
		},
		CommentText: commentEntity.Content,
		CreateDate:  commentEntity.CreatedAt,
	}, err
}

func (commentCtl *CommentCtl) QueryListByVideoId(videoId int64) ([]model.Comment, error) {

	// 数据库实体
	var commentEntitis []Comment

	// 查库
	err := dbProvider.GetDB().Preload("User").Where("video_id = ?", videoId).Find(&commentEntitis).Error

	// model实体
	comments := make([]model.Comment, len(commentEntitis))
	for i := 0; i < len(commentEntitis); i++ {
		comment := commentEntitis[i]
		comments[i] = model.Comment{
			ID: comment.ID,
			User: model.User{
				ID:            comment.User.ID,
				Name:          comment.User.UserName,
				FollowCount:   comment.User.FollowCount,
				FollowerCount: comment.User.FollowerCount,
			},
			CommentText: comment.Content,
			CreateDate:  comment.CreatedAt,
		}
	}

	return comments, err
}
