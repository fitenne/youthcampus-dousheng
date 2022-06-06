package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fitenne/youthcampus-dousheng/internal/service"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
)

type FavoriteRequest struct {
	VideoId    int64 `form:"video_id"`
	ActionType int   `form:"action_type"`
}

func FavoriteAction(c *gin.Context) {
	var req FavoriteRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: http.StatusBadRequest,
			StatusMsg:  "BadRequest",
		})
		return
	}

	// userId 校验
	userId := c.GetInt64("userID")
	// action_type: 1-点赞，2-取消点赞 其余返回异常
	var serverErr error = nil
	switch req.ActionType {
	case 1: // 点赞
		// 创建点赞表实体
		newfavorite := model.Favorite{UserId: int64(userId), VideoId: int64(req.VideoId)}

		// 调用点赞接口
		// 如果不为重复，就执行点赞
		if !service.CheckRepeatFavorite(userId, req.VideoId, &newfavorite) {
			serverErr = service.CreateFavoriteAction(req.VideoId, &newfavorite)
			// 异常分支处理：操作异常
			if serverErr != nil {
				log.Println("FavoriteAction:  流程异常" + serverErr.Error())
				c.JSON(http.StatusOK, Response{
					StatusCode: 3,
					StatusMsg:  "error " + serverErr.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "点赞成功"})
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "重复点赞"})
		}

		// 返回结果
		return

	case 2: // 取消点赞
		// 创建点赞表实体
		newfavorite := model.Favorite{UserId: int64(userId), VideoId: req.VideoId}
		// 调用删除接口
		serverErr = service.DeleteFavoriteAction(userId, req.VideoId, &newfavorite)

		// 异常分支处理：操作异常
		if serverErr != nil {
			log.Println("FavoriteAction: delete 流程异常" + serverErr.Error())
			c.JSON(http.StatusOK, Response{
				StatusCode: 3,
				StatusMsg:  "error " + serverErr.Error(),
			})
			return
		}
		// 正常返回
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "success"})
		return

	default: // 异常分支处理，操作类型异常
		c.JSON(http.StatusOK, Response{
			StatusCode: 2,
			StatusMsg:  "Invalid param action_type",
		})
		return
	}
}

func FavoriteList(c *gin.Context) {

	// token := c.DefaultQuery("token", "")
	userIdQuery := c.DefaultQuery("user_id", "")

	// 参数处理
	userId, err := strconv.ParseInt(userIdQuery, 10, 64)
	if err != nil {
		log.Println("videoList: videoId 转换异常" + err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 2,
			StatusMsg:  "Invalid param <video_id>: " + userIdQuery,
		})
		return
	}

	// 获取点赞视频列表
	VideoListDTOs, err := service.FavoriteVideoList(userId)

	// 异常分支处理：操作异常
	if err != nil {
		log.Println("VideoList: 函数流程异常" + err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 3,
			StatusMsg:  "server error " + err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "点赞视频列表"},
		VideoList: VideoListDTOs,
	})

}
