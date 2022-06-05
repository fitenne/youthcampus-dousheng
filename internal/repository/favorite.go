package repository

import (
	// "fmt"

	"fmt"

	"gorm.io/gorm"

	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

// type FavoriteList struct {
// 	gorm.Model
// 	Userid  int64
// 	Videoid int64
// }

type FavoriteCtl struct{}

var favoriteCtl FavoriteCtl

func GetFavoriteCtl() model.FavoriteCtl {
	return &favoriteCtl
}

// 创建一条点赞
func (favoriteCtl *FavoriteCtl) CreateFavoriteAction(videoId int64, newfavorite *model.Favorite) error {

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
	if err := tx.Where("user_id=? and video_id=?", userId, videoId).Delete(&newfavorite).Error; err != nil {
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
	dbProvider.GetDB().Model(&newfavorite).Where("user_id = ? and video_id = ?", userId, videoId).Count(&count)

	return count > 0
}
func (favoriteCtl *FavoriteCtl) FavoriteVideoList(userId int64) ([]model.Video, error) {
	var videolistEntitis []model.Video

	err := dbProvider.GetDB().Preload("Author").Joins("JOIN favorites ON favorites.video_id = videos.id AND favorites.user_id = ?", userId).Find(&videolistEntitis).Error

	for i := 0; i < len(videolistEntitis); i++ {

		isFollow, err := GetDealerFollow().CheckHasFollowed(userId, videolistEntitis[i].Author.ID)
		if isFollow {
			videolistEntitis[i].Author.IsFollow = isFollow
		} else {
			videolistEntitis[i].Author.IsFollow = false
		}
		if err != nil {
			return nil, err
		}
	}
	return videolistEntitis, err
}

// 根据结构体创建表
func (favoriteCtl *FavoriteCtl) CreateTableTest() error {
	err := dbProvider.GetDB().AutoMigrate(&model.Favorite{})
	if err != nil {
		fmt.Println(err)
	}
	return err
}
