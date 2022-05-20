package model

import (
	"time"
)

type Comment struct {
	ID          int64     `json:"id,omitempty"`
	User        User      `json:"user"`
	VideoId     int64     `json:"video_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"create_date"`
}

// 对数据库的修改
type CommentCtl interface {
	Publish(comment *Comment) error
	DeleteById(commentId int64) error
	QueryListByVideoId(videoId int64) ([]Comment, error)
}
