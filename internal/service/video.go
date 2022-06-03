package service

import (
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

func VideoPublish(c *gin.Context, data *multipart.FileHeader, playUrl string, authorID int64) (int64,  error ){

	// 存储视频文件
	if err := c.SaveUploadedFile(data, playUrl); err != nil {
		return -1, err
	}

	// 写入数据库
	video := &model.Video{
		AuthorID: authorID,
		PlayUrl:  playUrl,
		CoverUrl:      "",
		// FavoriteCount: 0,
		// CommentCount:  0,
		// CreatedAt:     0,
		// DeletedAt    : ,

	}

	videoId, err := repository.GetVideoCtl().Create(video)
	if err != nil {
		return -1, err
	}

	return videoId, nil
}
