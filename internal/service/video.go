package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"time"

	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
)

var StroageHost string

func PublishVideo(c *gin.Context, data *multipart.FileHeader, title string, authorID int64) (int64, error) {
	m := md5.Sum([]byte(fmt.Sprint(time.Now().UnixMicro(), data.Filename)))
	saveTo := hex.EncodeToString(m[:]) + filepath.Ext(data.Filename)

	// 存储视频文件
	if err := c.SaveUploadedFile(data, filepath.Join("./public", saveTo)); err != nil {
		return -1, err
	}

	playUrl := url.URL{
		Scheme:      "http",
		Host:        StroageHost,
		Path:        filepath.Join("/static", saveTo),
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	// video 信息写入数据库
	video := &model.Video{
		Title:    title,
		AuthorID: authorID,
		PlayUrl:  playUrl.String(),
		CoverUrl: "",
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

func GetVideos(authorID int64) ([]*model.Video, error) {
	videos, err := repository.GetVideoCtl().GetVideoByAuthorId(int(authorID))
	if err != nil {
		return nil, err
	}
	return videos, nil
}
