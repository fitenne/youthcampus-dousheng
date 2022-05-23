package controller_test

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/fitenne/youthcampus-dousheng/internal/common"
	"github.com/fitenne/youthcampus-dousheng/internal/common/jwt"
	"github.com/fitenne/youthcampus-dousheng/internal/common/mid"
	"github.com/fitenne/youthcampus-dousheng/internal/controller"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/internal/service/mocks"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var goodUser = repository.User{
	ID:            1,
	UserName:      "alice",
	Salt:          []byte("10086"),
	Password:      []byte("10086"),
	FollowCount:   0,
	FollowerCount: 0,
}
var notExistsUser = repository.User{
	ID:       1,
	UserName: "unknown",
	Salt:     []byte("10086"),
	Password: []byte("10086"),
}
var tooLongNameUser = repository.User{
	UserName: "123456789012345678901234567890123",
}

func setupRouter() *gin.Engine {
	g := gin.Default()
	router := g.Group("/douyin")

	router.GET("/user", mid.JWTAuthMiddleware(), controller.UserInfo)
	router.POST("/user/register", controller.Register)
	router.POST("/user/login", controller.Login)

	return g
}

func mock(t *testing.T) *mocks.MockUserCtl {
	mockCtrl := gomock.NewController(t)
	mockUserCtl := mocks.NewMockUserCtl(mockCtrl)

	notExistsName := gomock.Eq(notExistsUser.UserName)
	lenZero := gomock.Len(0)
	mockUserCtl.EXPECT().
		Create(notExistsName, lenZero, gomock.Any()).AnyTimes().
		Return(int64(0), errors.New("empty pass"))
	mockUserCtl.EXPECT().
		Create(notExistsName, gomock.Any(), lenZero).AnyTimes().
		Return(int64(0), errors.New("empty salt"))
	mockUserCtl.EXPECT().
		Create(gomock.Not(notExistsName), gomock.Any(), gomock.Any()).AnyTimes().
		Return(int64(0), errors.New("username already exists"))
	mockUserCtl.EXPECT().
		Create(notExistsName, gomock.Not(lenZero), gomock.Not(lenZero)).AnyTimes().
		Return(int64(1), nil)

	goodName := gomock.Not(gomock.Eq(notExistsUser.UserName))
	goodID := gomock.Not(gomock.Eq(notExistsUser.ID))
	mockUserCtl.EXPECT().
		QueryByID(gomock.Not(goodID)).AnyTimes().
		Return(model.User{}, common.ErrUserNotExists{})
	mockUserCtl.EXPECT().
		QueryByID(goodID).AnyTimes().
		Return(model.User{}, nil)

	mockUserCtl.EXPECT().
		QueryByName(gomock.Not(goodName)).AnyTimes().
		Return(model.User{}, common.ErrUserNotExists{})
	mockUserCtl.EXPECT().
		QueryByName(goodName).AnyTimes().
		Return(model.User{}, nil)

	mockUserCtl.EXPECT().
		QueryCredentialsByName(gomock.Not(goodName)).AnyTimes().
		Return(int64(0), nil, nil, common.ErrUserNotExists{})
	mockUserCtl.EXPECT().
		QueryCredentialsByName(goodName).AnyTimes().
		Return(int64(0), goodUser.Password, goodUser.Salt, nil)

	return mockUserCtl
}

func TestRegister(t *testing.T) {
	g := setupRouter()
	mockUserCtl := mock(t)
	repository.GetUserCtl = func() model.UserCtl {
		return mockUserCtl
	}

	parseBody := func(r io.Reader, t *testing.T) controller.UserLoginResponse {
		body, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatal(err.Error())
		}
		resp := controller.UserLoginResponse{}
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatal(err.Error())
		}
		return resp
	}

	t.Run("normal register", func(t *testing.T) {
		username := notExistsUser.UserName
		password := notExistsUser.Password
		url := url.URL{Path: "/douyin/user/register"}
		query := url.Query()
		query.Add("username", username)
		query.Add("password", string(password))
		url.RawQuery = query.Encode()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", url.String(), nil)
		g.ServeHTTP(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("/douyin/user/register: want status code: %v, but got %v", http.StatusOK, w.Result().StatusCode)
		}

		resp := parseBody(w.Body, t)
		if resp.UserId != notExistsUser.ID {
			t.Errorf("/douyin/user/register: want id: %v, but got %v", notExistsUser.ID, resp.UserId)
		}
		token, err := jwt.GenToken(resp.UserId)
		if err != nil {
			t.Fatal(err.Error())
		}
		if resp.Token != token {
			t.Errorf("/douyin/user/register: want token: %v, but got %v", token, resp.Token)
		}
	})

	type Args struct {
		username, password string
	}
	type Case struct {
		name     string
		args     Args
		wantCode int
	}
	badRequestCases := []Case{
		{
			name: "too long name",
			args: Args{
				username: tooLongNameUser.UserName,
				password: "10086",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "too long pass",
			args: Args{
				username: "12345678901234567890123456789012",
				password: "123456789012345678901234567890123",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "empty name",
			args: Args{
				username: "",
				password: "10086",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "empty pass",
			args: Args{
				username: "10086",
				password: "",
			},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range badRequestCases {
		t.Run(tt.name, func(t *testing.T) {
			url := url.URL{Path: "/douyin/user/register"}
			query := url.Query()
			query.Add("username", tt.args.username)
			query.Add("password", tt.args.password)
			url.RawQuery = query.Encode()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", url.String(), nil)
			g.ServeHTTP(w, req)

			if w.Result().StatusCode != http.StatusOK {
				t.Errorf("/douyin/user/register: want status code: %v, but got %v", http.StatusOK, w.Result().StatusCode)
			}

			resp := parseBody(w.Body, t)
			if resp.Response.StatusCode != http.StatusBadRequest {
				t.Errorf("/douyin/user/register: want response code: %v, but got %v", http.StatusBadRequest, resp.Response.StatusCode)
			}
		})
	}
}
