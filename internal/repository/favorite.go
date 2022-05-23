package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

type FavoriteList struct {
	gorm.Model
	Userid  int64
	Videoid int64
}

type FavoriteCtl struct{}

var favoriteCtl FavoriteCtl

func GetFavoriteCtl() model.FavoriteCtl {
	return &favoriteCtl
}

// 创建一条点赞
func (favoriteCtl *FavoriteCtl) CreateFavoriteAction(videoId int64, newfavorite *model.Favorite) error {
	err := dbProvider.GetDB().AutoMigrate(&FavoriteList{})
	if err != nil {
		fmt.Println(err)
	}
	// 开启事务
	tx := dbProvider.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	// 点赞关系表中创建一条记录
	if err := tx.Create(&newfavorite).Error; err != nil {
		tx.Rollback()
		return err
	}
	var video model.Video
	// Video表中点赞数+1
	if err := tx.Model(&video).Where("id=?", videoId).Update("favorite_count", gorm.Expr("favorite_count  + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

// 删除一条点赞(软删除)
func (favoriteCtl *FavoriteCtl) DeleteFavoriteAction(userId int64, videoId int64, newfavorite *model.Favorite) error {

	// 开启事务
	tx := dbProvider.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	// 点赞关系表中删除一条记录
	if err := tx.Where("userid=? and videoid=?", userId, videoId).Delete(&newfavorite).Error; err != nil {
		tx.Rollback()
		return err
	}
	var video model.Video
	// Video表中点赞数-1
	if err := tx.Model(&video).Where("id=?", videoId).Update("favorite_count", gorm.Expr("favorite_count  + ?", -1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 查询是否已经点赞
func (favoriteCtl *FavoriteCtl) CheckRepeatFavorite(userId int64, videoId int64, newfavorite *model.Favorite) bool {

	var count int64
	dbProvider.GetDB().Model(&newfavorite).Where("userid = ? and videoid = ?", userId, videoId).Count(&count)

	return count > 0
}
