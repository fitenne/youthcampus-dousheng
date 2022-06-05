package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/fitenne/youthcampus-dousheng/internal/common"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/internal/service"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User model.User `json:"user"`
}

type RegisterRequest struct {
	Username string `form:"username" binding:"required,gt=0,lte=32"`
	Password string `form:"password" binding:"required,gt=0,lte=32"`
}

type LoginRequest struct {
	Username string `form:"username" binding:"required,gt=0,lte=32"`
	Password string `form:"password" binding:"required,gt=0,lte=32"`
}

func Register(c *gin.Context) {
	var regReq RegisterRequest
	if err := c.ShouldBindQuery(&regReq); err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusBadRequest,
				StatusMsg:  "invalid request",
			},
		})
		return
	}

	exists, err := service.UserExists(regReq.Username)
	if exists {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusForbidden,
				StatusMsg:  "User already exist",
			},
		})
		return
	}
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusInternalServerError,
				StatusMsg:  "internal server error",
			},
		})
		return
	}

	id, token, err := service.UserRegister(regReq.Username, regReq.Password)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusInternalServerError,
				StatusMsg:  "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		UserId: id,
		Token:  token,
	})
}

func Login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindQuery(&loginReq); err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusBadRequest,
				StatusMsg:  "invalid request",
			},
		})
		return
	}

	id, token, err := service.UserLogin(loginReq.Username, loginReq.Password)
	if err != nil {
		if errors.Is(err, common.ErrUserNotExists{}) {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: http.StatusUnauthorized,
					StatusMsg:  "username password mismatch",
				},
			})
		}

		log.Print(err.Error())
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusInternalServerError,
				StatusMsg:  "internal server error",
			},
		})
		return
	}
	if token == "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusUnauthorized,
				StatusMsg:  "invalid credentials",
			},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   id,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	rid := c.Query("user_id")
	id, err := strconv.Atoi(rid)
	if rid == "" || err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusBadRequest,
				StatusMsg:  "invalid credentials",
			},
		})
		return
	}

	u, err := repository.GetUserCtl().QueryByID(int64(id))
	//! not implemented, is_follow = GetCtl().QueryFollow(me, id)
	if err != nil {
		if errors.Is(err, common.ErrUserNotExists{}) {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: http.StatusBadRequest,
					StatusMsg:  "invalid credentials",
				},
			})
			return
		}

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusInternalServerError,
				StatusMsg:  "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     u,
	})
}
