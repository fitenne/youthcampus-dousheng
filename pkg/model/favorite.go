package model

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	Userid  int64
	Videoid int64
}

// 对数据库的修改
type FavoriteCtl interface {
	CreateFavoriteAction(videoId int64, newfavorite *Favorite) error
	DeleteFavoriteAction(userId int64, videoId int64, newfavorite *Favorite) error
	CheckRepeatFavorite(userId int64, videoId int64, newfavorite *Favorite) bool
}
