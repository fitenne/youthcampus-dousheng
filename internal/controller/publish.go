package controller

import (
	"fmt"
	"github.com/fitenne/youthcampus-dousheng/internal/common/jwt"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
)

type VideoListResponse struct {
	Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	// 判断用户是否存在
	token := c.PostForm("token")
	//if _, exist := usersLoginInfo[token]; !exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//	return
	//}
	myClaims, exist := jwt.ParseToken(token)
	if exist != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	//
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	log.Println("publish func, filename:", filename)
	// user := usersLoginInfo[token]
	// finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	finalName := fmt.Sprintf("%d_%s", myClaims.UserID, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	log.Println(finalName + " uploaded successfully")
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
