package controller

import (
	// "fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fitenne/youthcampus-dousheng/internal/service"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
)

func FavoriteAction(c *gin.Context) {

	userIdQuery := c.DefaultQuery("user_id", "")
	videoIdQuery := c.DefaultQuery("video_id", "")
	actionType := c.DefaultQuery("action_type", "")

	// userId 校验
	userId, err := strconv.ParseInt(userIdQuery, 10, 64)
	if err != nil {
		log.Println("FavoriteAction: userId 转换异常" + err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 2,
			StatusMsg:  "Invalid param <user_id>: " + userIdQuery,
		})
		return
	}

	// action_type: 1-点赞，2-取消点赞 其余返回异常
	var serverErr error = nil
	switch actionType {
	case "1": // 点赞

		// videoId 校验
		videoId, err := strconv.ParseInt(videoIdQuery, 10, 64)
		if err != nil {
			log.Println("FavoriteAction: videoId 转换异常" + err.Error())
			c.JSON(http.StatusOK, Response{
				StatusCode: 2,
				StatusMsg:  "Invalid param <comment_id>: " + videoIdQuery + err.Error(),
			})
			return
		}

		// 创建点赞表实体
		newfavorite := model.Favorite{UserId: int64(userId), VideoId: int64(videoId)}

		// 调用点赞接口
		// 如果不为重复，就执行点赞
		if !service.CheckRepeatFavorite(userId, videoId, &newfavorite) {
			serverErr = service.CreateFavoriteAction(videoId, &newfavorite)
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

	case "2": // 取消点赞

		// videoId 校验
		videoId, err := strconv.ParseInt(videoIdQuery, 10, 64)
		if err != nil {
			log.Println("FavoriteAction: videoId 转换异常" + err.Error())
			c.JSON(http.StatusOK, Response{
				StatusCode: 2,
				StatusMsg:  "Invalid param <comment_id>: " + videoIdQuery + err.Error(),
			})
			return
		}

		// 创建点赞表实体
		newfavorite := model.Favorite{UserId: int64(userId), VideoId: int64(videoId)}
		// 调用删除接口
		serverErr = service.DeleteFavoriteAction(userId, videoId, &newfavorite)

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
	// case "3":
	// 	// 创建follow表
	// 	err := service.CreateTableTest()
	// 	if err != nil {
	// 		c.JSON(http.StatusOK, Response{
	// 			StatusCode: 3,
	// 			StatusMsg:  "error " + serverErr.Error(),
	// 		})
	// 	}
	default: // 异常分支处理，操作类型异常
		c.JSON(http.StatusOK, Response{
			StatusCode: 2,
			StatusMsg:  "Invalid param <action_type>: " + actionType,
		})
		return
	}
}

// FavoriteList all users have same favorite video list
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
