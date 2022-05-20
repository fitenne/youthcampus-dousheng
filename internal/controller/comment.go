package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fitenne/youthcampus-dousheng/internal/service"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	// 参数获取
	token := c.DefaultQuery("token", "")
	userIdQuery := c.DefaultQuery("user_id", "")
	videoIdQuery := c.DefaultQuery("video_id", "")
	actionType := c.DefaultQuery("action_type", "")

	// 形成参数表
	params := map[string]string{
		"token":       token,
		"user_id":     userIdQuery,
		"video_id":    videoIdQuery,
		"action_type": actionType,
	}

	// 参数校验
	for key, value := range params {
		if value == "" {
			c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "Param " + key + " can't be empty"})
			return
		}
	}

	// 检查token
	if _, exist := usersLoginInfo[params["token"]]; exist {

		// action_type: 1-发布评论，2-删除评论 其余返回异常
		var err error
		switch actionType {
		case "1": // 发布评论

			// userId 校验
			userId, err := strconv.ParseInt(userIdQuery, 10, 64)
			if err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 2,
					StatusMsg:  "Invalid param <comment_id>: " + userIdQuery,
				})
				return
			}

			// videoId 校验
			videoId, err := strconv.ParseInt(videoIdQuery, 10, 64)
			if err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 2,
					StatusMsg:  "Invalid param <comment_id>: " + videoIdQuery,
				})
				return
			}

			// 获取评论内容
			commentText := c.Query("comment_text")

			// 调用发布接口
			err = service.Publish(videoId, &model.Comment{
				User: model.User{
					ID: userId,
				},
				CommentText: commentText,
				CreatedAt:   time.Now(),
			})

		case "2": // 删除评论

			// 获取评论id
			commentIdQuery := c.Query("comment_id")

			// 参数处理
			commentId, err := strconv.ParseInt(commentIdQuery, 10, 64)
			if err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 2,
					StatusMsg:  "Invalid param <comment_id>: " + commentIdQuery,
				})
				return
			}

			// 调用删除接口
			err = service.DeleteById(commentId)

		default: // 异常分支处理，操作类型异常
			c.JSON(http.StatusOK, Response{
				StatusCode: 2,
				StatusMsg:  "Invalid param <action_type>: " + actionType,
			})
			return
		}

		// 异常分支处理：操作异常
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 3,
				StatusMsg:  "server error " + err.Error(),
			})
			return
		}

		// 正常返回
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	// 获取评论id
	videoIdQuery := c.DefaultQuery("video_id", "")

	// 参数处理
	videoId, err := strconv.ParseInt(videoIdQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 2,
			StatusMsg:  "Invalid param <video_id>: " + videoIdQuery,
		})
		return
	}

	// 获取评论 替换DemoComments
	commentDTOs, err := service.QueryListByVideoId(videoId)

	// 异常分支处理：操作异常
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 3,
			StatusMsg:  "server error " + err.Error(),
		})
		return
	}

	// 类型转换
	comments := make([]Comment, len(commentDTOs))

	// 返回结果
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
