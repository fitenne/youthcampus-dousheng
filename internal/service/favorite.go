package service

import (
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

var favoriteCtl = repository.GetFavoriteCtl()

func CreateFavoriteAction(videoId int64, newfavorite *model.Favorite) error {
	return favoriteCtl.CreateFavoriteAction(videoId, newfavorite)
}
func DeleteFavoriteAction(userId int64, videoId int64, newfavorite *model.Favorite) error {
	return favoriteCtl.DeleteFavoriteAction(userId, videoId, newfavorite)
}
func CheckRepeatFavorite(userId int64, videoId int64, newfavorite *model.Favorite) bool {
	return favoriteCtl.CheckRepeatFavorite(userId, videoId, newfavorite)
}
