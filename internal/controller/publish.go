package controller

import (
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/fitenne/youthcampus-dousheng/internal/service"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []model.Video `json:"video_list"`
}

type PublishRequest struct {
	Data  multipart.FileHeader `form:"data"`
	Title string               `form:"title"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	var req PublishRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Panicln(err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: http.StatusBadRequest,
			StatusMsg:  "BadRequest",
		})
	}

	// 存入数据库
	// filename := filepath.Base(data.Filename)
	// log.Println("publish func, filename:", filename)
	// // user := usersLoginInfo[token]
	// // finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	// finalName := fmt.Sprintf("%d_%s", userID, filename)
	userID := c.GetInt64("userID")
	videoID, err := service.PublishVideo(c, &req.Data, req.Title, userID)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 上传成功，返回结果
	log.Println(videoID, req.Data.Filename+" uploaded successfully")
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  req.Data.Filename + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {

	//// 判断用户是否存在
	//token := c.PostForm("token")
	//myClaims, exist := jwt.ParseToken(token)
	//if exist != nil {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//	return
	//}
	//userID := myClaims.UserID

	userID, _ := strconv.Atoi(c.Query("user_id"))
	// 查询用户发布视频列表
	videos, err := service.GetVideos(int64(userID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
			},
			VideoList: nil,
		})
		return
	}

	// 查询成功，返回结果
	var videosList []model.Video
	for _, v := range videos {
		newVideo := model.Video{
			ID:            v.ID,
			AuthorID:      v.AuthorID,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
		}
		videosList = append(videosList, newVideo)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videosList,
	})
}
