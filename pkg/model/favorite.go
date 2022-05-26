package model

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserId int64
	// Userid  int64
	VideoId int64
}

// 对数据库的修改
type FavoriteCtl interface {
	CreateFavoriteAction(videoId int64, newfavorite *Favorite) error
	DeleteFavoriteAction(userId int64, videoId int64, newfavorite *Favorite) error
	CheckRepeatFavorite(userId int64, videoId int64, newfavorite *Favorite) bool
	CreateTableTest() error
	FavoriteVideoList(userId int64) ([]Video, error)
}
