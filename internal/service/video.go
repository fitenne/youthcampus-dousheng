package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
)

var StroageHost string

// input public/filename.ext, output public/img/filename.png
func generateThumbnailFromVideo(filename, ext string) {
	// this just works
	var ffmpeg = "/usr/bin/ffmpeg"
	c := exec.Command(ffmpeg, "-i", fmt.Sprintf("public/%v%v", filename, ext), "-vf", "thumbnail", "-frames:v", "1", fmt.Sprintf("public/img/%v.png", filename))
	if err := c.Run(); err != nil {
		log.Println(err.Error())
	}
}

func PublishVideo(c *gin.Context, data *multipart.FileHeader, title string, authorID int64) (int64, error) {
	m := md5.Sum([]byte(fmt.Sprint(time.Now().UnixMicro(), data.Filename)))
	filename := hex.EncodeToString(m[:])
	saveTo := filename + filepath.Ext(data.Filename)

	// 存储视频文件
	if err := c.SaveUploadedFile(data, filepath.Join("./public", saveTo)); err != nil {
		return -1, err
	}
	generateThumbnailFromVideo(filename, filepath.Ext(data.Filename))

	playUrl := url.URL{
		Scheme: "http",
		Host:   StroageHost,
		Path:   filepath.Join("/static", saveTo),
	}
	// this just works
	coverUrl := url.URL{
		Scheme: "http",
		Host:   StroageHost,
		Path:   filepath.Join("/static/img", filename+".png"),
	}

	// video 信息写入数据库
	video := &model.Video{
		Title:    title,
		AuthorID: authorID,
		PlayUrl:  playUrl.String(),
		CoverUrl: coverUrl.String(),
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
