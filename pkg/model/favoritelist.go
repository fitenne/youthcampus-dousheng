package model

import (
	"gorm.io/gorm"
)

type FavoriteList struct {
	gorm.Model
	Userid  int64
	Videoid int64
}

// 对数据库的修改
type FavoriteCtl interface {
}
