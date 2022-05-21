package controller

import (
	"log"
	"net/http"

	"github.com/fitenne/youthcampus-dousheng/internal/common/mid"
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

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if username == "" || password == "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusBadRequest,
				StatusMsg:  "invalid request",
			},
		})
		return
	}

	exists, err := service.UserExists(username)
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

	id, token, err := service.UserRegister(username, password)
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
	username := c.Query("username")
	password := c.Query("password")

	if username == "" || password == "" {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: http.StatusBadRequest,
				StatusMsg:  "invalid request",
			},
		})
		return
	}

	id, token, err := service.UserLogin(username, password)
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
	id := c.GetInt64(mid.UserIDKey)

	u, err := repository.GetUserCtl().QueryByID(id)
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

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     u,
	})
}
