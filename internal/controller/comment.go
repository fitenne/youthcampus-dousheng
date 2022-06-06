package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fitenne/youthcampus-dousheng/internal/common/jwt"

	"github.com/fitenne/youthcampus-dousheng/internal/service"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
)

// CommentActionRequest 增删评论请求体
type CommentActionRequest struct {
	Token      string `form:"token" json:"token" binding:"required"`
	VideoId    int64  `form:"video_id" json:"video_id" binding:"required"`
	ActionType int64  `form:"action_type" json:"action_type" binding:"required"`
}

// CommentActionResponse 增删评论响应体
type CommentActionResponse struct {
	Response
	Comment model.Comment `json:"comment,omitempty"`
}

// CommentListResponse 列表查询响应体
type CommentListResponse struct {
	Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

// CommentAction 评论增删
func CommentAction(c *gin.Context) {

	// 获取参数
	var comActReq CommentActionRequest
	if err := c.ShouldBindQuery(&comActReq); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "miss params"})
		return
	}

	// 检查token
	claims, err := jwt.ParseToken(comActReq.Token)
	userId := claims.UserID
	if err != nil {
		log.Println("controller.CommentAction|token parse error|jwt.ParseToken: " + err.Error())
		c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "invalid token"})
		return
	}

	// action_type: 1-发布评论，2-删除评论 其余返回异常
	switch comActReq.ActionType {
	case 1: // 发布评论

		// 生成comment
		comment := model.Comment{
			VideoId:    comActReq.VideoId,
			UserID:     userId,
			Content:    c.DefaultQuery("comment_text", ""),
			CreateDate: time.Now().Format("01-02"),
		}

		// 调用发布接口, 异常处理
		if serverErr := service.Publish(userId, &comment); serverErr != nil {
			log.Println("controller.CommentAction|server error|service.Publish: " + serverErr.Error())
			c.JSON(http.StatusOK, Response{StatusCode: 3, StatusMsg: "server error: " + serverErr.Error()})
			return
		}

		// 返回结果
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "success"},
			Comment:  comment,
		})
		return

	case 2: // 删除评论

		// 获取评论id
		commentIdQuery, ok := c.GetQuery("comment_id")
		if !ok {
			c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "invalid comment_id"})
			return
		}

		// 转换commentId
		commentId, err := strconv.ParseInt(commentIdQuery, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "invalid comment_id"})
			return
		}

		// 调用删除接口, 异常处理
		if serverErr := service.DeleteById(commentId, userId, comActReq.VideoId); serverErr != nil {
			log.Println("controller.CommentAction|server error|service.DeleteById: " + serverErr.Error())
			c.JSON(http.StatusOK, Response{StatusCode: 3, StatusMsg: "server error: " + serverErr.Error()})
			return
		}

		// 正常返回
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "success"})
		return

	default: // 异常分支处理，操作类型异常
		c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "invalid action_type"})
		return
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	// 检查token
	claims, err := jwt.ParseToken(c.DefaultQuery("token", ""))
	userId := claims.UserID
	if err != nil {
		log.Println("controller.CommentList|token parse error|jwt.ParseToken: " + err.Error())
		c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "invalid token"})
		return
	}

	// 获取评论id
	videoIdQuery, ok := c.GetQuery("video_id")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "invalid video_id"})
		return
	}

	// video转换
	videoId, err := strconv.ParseInt(videoIdQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "invalid video_id"})
		return
	}

	// 获取评论, 异常处理
	comments, err := service.QueryListByVideoId(videoId, userId)
	if err != nil {
		log.Println("controller.CommentList|server error|service.QueryListByVideoId: " + err.Error())
		c.JSON(http.StatusOK, Response{StatusCode: 3, StatusMsg: "server error: " + err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: comments,
	})
}
