package service

import (
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

// return LatestTime for next query
func GetFeedUntil(latestTime int64, topN int) (videoList []model.Video, err error) {
	videosDto, err := repository.GetVideoCtl().GetVideoList(latestTime, topN)

	videoList = make([]model.Video, 0, len(videosDto))
	for _, v := range videosDto {
		videoList = append(videoList, *v)
	}
	return videoList, err
}
