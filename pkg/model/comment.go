package model

import (
	"time"
)

type Comment struct {
	ID          int64     `json:"id,omitempty"`
	User        User      `json:"user"`
	CommentText string    `json:"comment_text"`
	CreateDate  time.Time `json:"create_date"`
}

// 对数据库的修改
type CommentCtl interface {
	Publish(videoId int64, comment *Comment) error
	DeleteById(commentId int64) error
	QueryById(commentId int64) (*Comment, error)																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																
	QueryListByVideoId(videoId int64) ([]Comment, error)
}																																																																																																																																																																							
