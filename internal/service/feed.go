package service

import (
	"errors"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

func Feed(user_id, latest int64) ([]*model.Video, int64, error) {
	videos, err := repository.GetVideoCtl().GetVideoList(latest, 30)
	if err != nil {
		return nil, 0, err
	}
	if len(videos) == 0 {
		return nil, 0, errors.New("无video")
	}
	next_time := videos[len(videos)-1].CreatedAt
	//给videos的IsFavorite赋值
	for _, v := range videos {
		v.IsFavorite = repository.GetFavoriteCtl().CheckRepeatFavorite(user_id, v.ID, &model.Favorite{})
	}
	return videos, int64(next_time), nil
}
