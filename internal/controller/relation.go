package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fitenne/youthcampus-dousheng/pkg/model"

	"github.com/fitenne/youthcampus-dousheng/internal/common/code"

	"github.com/fitenne/youthcampus-dousheng/internal/service"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []model.User `json:"user_list"`
}

type RelationActionRequest struct {
	ToUserID   int64 `form:"to_user_id" binding:"required"`
	ActionType int64 `form:"action_type" binding:"required"`
}

func StatusResponse(c *gin.Context, code int32, msg string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: code,
		StatusMsg:  msg,
	})
}

// RelationAction 关注或取消关注操作
func RelationAction(c *gin.Context) {
	var req RelationActionRequest
	userID := c.GetInt64("userID")
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: http.StatusBadRequest,
			StatusMsg:  "BadRequest",
		})
		return
	}

	if req.ActionType == 1 { //关注
		err := service.FollowUser(userID, req.ToUserID)
		if err != nil {
			StatusResponse(c, -1, err.Error())
			return
		}
		StatusResponse(c, 0, code.Success.Msg())
		return
	} else if req.ActionType == 2 { //取消关注
		err := service.CancelFollowUser(userID, req.ToUserID)
		if err != nil {
			StatusResponse(c, -1, err.Error())
			return
		}
		StatusResponse(c, 0, code.Success.Msg())
		return
	} else {
		StatusResponse(c, -1, code.InvalidParameter.Msg())
		return
	}
}

// FollowList 获取关注列表
func FollowList(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		StatusResponse(c, -1, code.InvalidParameter.Msg())
		return
	}
	userList, err := service.GetFollowList(userID)
	if err != nil {
		StatusResponse(c, -1, err.Error())
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  code.Success.Msg(),
		},
		UserList: *userList,
	})
}

// FollowerList 获取粉丝列表
func FollowerList(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		StatusResponse(c, -1, code.InvalidParameter.Msg())
		return
	}
	userList, err := service.GetFollowerList(userID)
	if err != nil {
		StatusResponse(c, -1, err.Error())
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  code.Success.Msg(),
		},
		UserList: *userList,
	})
}
