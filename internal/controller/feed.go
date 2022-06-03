package controller

import (
	"github.com/fitenne/youthcampus-dousheng/internal/service"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []*model.Video `json:"video_list,omitempty"`
	NextTime  int64          `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	rid := c.Query("user_id")
	latestTime, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	if err != nil {
		latestTime = time.Now().Unix()
	}
	var videos []*model.Video
	var next_time int64
	if rid == "" {
		videoList, nextTime, err := service.Feed(-1, latestTime)
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				Response: Response{
					StatusCode: http.StatusInternalServerError,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		videos = videoList
		next_time = nextTime
	} else {
		userId, err := strconv.ParseInt(rid, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				Response: Response{
					StatusCode: http.StatusBadRequest,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		videoList, nextTime, err := service.Feed(userId, latestTime)
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				Response: Response{
					StatusCode: http.StatusInternalServerError,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		videos = videoList
		next_time = nextTime
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  next_time,
	})
}
