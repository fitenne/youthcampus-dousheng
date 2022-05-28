package model

import (
	"gorm.io/gorm"
)

type UserEntity struct {
	ID            int64  `json:"id,omitempty" gorm:"primaryKey"`
	UserName      string `json:"name,omitempty" gorm:"user_name"`
	FollowCount   int64  `json:"follow_count,omitempty" gorm:"follow_count"`
	FollowerCount int64  `json:"follower_count,omitempty" gorm:"follower_count"`
	IsFollow      bool   `json:"is_follow,omitempty" gorm:"-"`
}

func (UserEntity) TableName() string {
	return "users"
}

type Comment struct {
	ID         int64           `json:"id,omitempty" gorm:"primaryKey;comment:评论ID;autoIncrement;unique_index:create_time_index"`
	Content    string          `json:"content" gorm:"content;comment:评论内容;unique_index:create_time_index;not null"`
	UserID     int64           `json:"-" gorm:"user_id;comment:发布者ID;unique_index:create_time_index;not null"`
	VideoId    int64           `json:"video_id" gorm:"video_id;comment:视频ID;unique_index:create_time_index;not null"`
	CreateDate string          `json:"create_date" gorm:"create_date;comment:评论时间;unique_index:create_time_index;not null"`
	DeletedAt  *gorm.DeletedAt `json:"-" gorm:"index;comment:删除标记位;unique_index:create_time_index"`

	// 发布者
	User *UserEntity `json:"user" gorm:"ForeignKey:UserID"`
}

// CommentCtl 对数据库的修改
type CommentCtl interface {

	// Publish 发布接口
	Publish(comment *Comment) error

	// DeleteById 删除接口
	DeleteById(commentId int64) error
	QueryById(commentId int64) (*Comment, error)
	QueryListByVideoId(videoId int64) ([]Comment, error)
}
