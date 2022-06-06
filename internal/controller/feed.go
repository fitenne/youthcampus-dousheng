package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/fitenne/youthcampus-dousheng/internal/service"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
)

const (
	numVideoInResponse = 30
)

type FeedResponse struct {
	Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

type FeedRequest struct {
	LatestTime int64  `form:"latest_time"`
	Token      string `form:"token"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var req FeedRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{
				StatusCode: http.StatusBadRequest,
				StatusMsg:  "BadRequest",
			},
		})
		return
	}

	if req.LatestTime == 0 {
		// no latest_time specified
		req.LatestTime = time.Now().Unix() + 1
	}

	videos, err := service.GetFeedUntil(req.LatestTime, numVideoInResponse)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{
				StatusCode: http.StatusInternalServerError,
				StatusMsg:  "InternalServerError",
			},
		})
		return
	}

	nextTime := time.Now().Unix()
	if len(videos) != 0 {
		nextTime = int64(videos[len(videos)-1].CreatedAt)
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  nextTime,
	})
}
