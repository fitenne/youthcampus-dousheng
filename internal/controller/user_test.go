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
	Password:      []byte{0x11, 0x46, 0x5, 0x19, 0x70, 0x5b, 0x29, 0x7e, 0x4, 0x2c, 0xd, 0xa8, 0x43, 0xe0, 0x54, 0xb5, 0x85, 0x4b, 0xa3, 0xc, 0xfc, 0x3, 0xff, 0x51, 0x18, 0xba, 0x79, 0xa9, 0xf9, 0xc2, 0xbb, 0x64}, // sha256("1008610086")
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

	router.GET("/user/", mid.JWTAuthMiddleware(), controller.UserInfo)
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

	goodName := gomock.Eq(goodUser.UserName)
	goodID := gomock.Eq(goodUser.ID)
	mockUserCtl.EXPECT().
		QueryByID(gomock.Not(goodID)).AnyTimes().
		Return(model.User{}, common.ErrUserNotExists{})
	mockUserCtl.EXPECT().
		QueryByID(gomock.Eq(goodUser.ID)).AnyTimes().
		Return(model.User{
			ID:   goodUser.ID,
			Name: goodUser.UserName,
		}, nil)

	mockUserCtl.EXPECT().
		QueryByName(gomock.Not(goodName)).AnyTimes().
		Return(model.User{}, common.ErrUserNotExists{})
	mockUserCtl.EXPECT().
		QueryByName(goodName).AnyTimes().
		Return(model.User{
			ID:   goodUser.ID,
			Name: goodUser.UserName,
		}, nil)

	mockUserCtl.EXPECT().
		QueryCredentialsByName(gomock.Not(goodName)).AnyTimes().
		Return(int64(0), nil, nil, common.ErrUserNotExists{})
	mockUserCtl.EXPECT().
		QueryCredentialsByName(goodName).AnyTimes().
		Return(goodUser.ID, goodUser.Password, goodUser.Salt, nil)

	return mockUserCtl
}

var parseUserLoginResponse = func(r io.Reader, t *testing.T) controller.UserLoginResponse {
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

var parseUserResponse = func(r io.Reader, t *testing.T) controller.UserResponse {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err.Error())
	}
	resp := controller.UserResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatal(err.Error())
	}
	return resp
}

func TestRegister(t *testing.T) {
	g := setupRouter()
	mockUserCtl := mock(t)
	repository.GetUserCtl = func() model.UserCtl {
		return mockUserCtl
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

		resp := parseUserLoginResponse(w.Body, t)
		if resp.StatusCode != 0 {
			t.Errorf("/douyin/user/register: want StatusCode: %v, but got %v", 0, resp.UserId)
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

			resp := parseUserLoginResponse(w.Body, t)
			if resp.Response.StatusCode != http.StatusBadRequest {
				t.Errorf("/douyin/user/register: want response code: %v, but got %v", http.StatusBadRequest, resp.Response.StatusCode)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	g := setupRouter()
	mockUserCtl := mock(t)
	repository.GetUserCtl = func() model.UserCtl {
		return mockUserCtl
	}

	t.Run("normal login", func(t *testing.T) {
		url := url.URL{Path: "/douyin/user/login"}
		query := url.Query()
		query.Add("username", goodUser.UserName)
		query.Add("password", "10086")
		url.RawQuery = query.Encode()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", url.String(), nil)
		g.ServeHTTP(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("/douyin/user/login: want status code: %v, but got %v", http.StatusOK, w.Result().StatusCode)
		}

		resp := parseUserLoginResponse(w.Body, t)
		if resp.Response.StatusCode != 0 {
			t.Errorf("/douyin/user/login: want response code: %v, but got %v", http.StatusOK, resp.Response.StatusCode)
		}
	})

	unAuth := []struct {
		name               string
		username, password string
	}{
		{"mismatch name pass 1", goodUser.UserName, "1234"},
		{"mismatch name pass 2", "bob", "10086"},
	}
	for _, tt := range unAuth {
		t.Run(tt.name, func(t *testing.T) {
			url := url.URL{Path: "/douyin/user/login"}
			query := url.Query()
			query.Add("username", tt.username)
			query.Add("password", tt.password)
			url.RawQuery = query.Encode()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", url.String(), nil)
			g.ServeHTTP(w, req)

			if w.Result().StatusCode != http.StatusOK {
				t.Errorf("/douyin/user/login: want status code: %v, but got %v", http.StatusOK, w.Result().StatusCode)
			}

			resp := parseUserLoginResponse(w.Body, t)
			if resp.Response.StatusCode != http.StatusUnauthorized {
				t.Errorf("/douyin/user/login: want response code: %v, but got %v", http.StatusUnauthorized, resp.Response.StatusCode)
			}
		})
	}

	badRequest := []struct {
		name               string
		username, password string
	}{
		{"too long name", tooLongNameUser.UserName, "10086"},
		{"too long pass", "alice", tooLongNameUser.UserName},
	}
	for _, tt := range badRequest {
		t.Run(tt.name, func(t *testing.T) {
			url := url.URL{Path: "/douyin/user/login"}
			query := url.Query()
			query.Add("username", tt.username)
			query.Add("password", tt.password)
			url.RawQuery = query.Encode()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", url.String(), nil)
			g.ServeHTTP(w, req)

			if w.Result().StatusCode != http.StatusOK {
				t.Errorf("/douyin/user/login: want status code: %v, but got %v", http.StatusOK, w.Result().StatusCode)
			}

			resp := parseUserLoginResponse(w.Body, t)
			if resp.Response.StatusCode != http.StatusBadRequest {
				t.Errorf("/douyin/user/login: want response code: %v, but got %v", http.StatusBadRequest, resp.Response.StatusCode)
			}
		})
	}
}

func TestUserInfo(t *testing.T) {
	g := setupRouter()
	mockUserCtl := mock(t)
	repository.GetUserCtl = func() model.UserCtl {
		return mockUserCtl
	}

	t.Run("normal info", func(t *testing.T) {
		token, err := jwt.GenToken(1)
		if err != nil {
			t.Fatal(err.Error())
		}

		url := url.URL{Path: "/douyin/user/"}
		query := url.Query()
		query.Add("token", token)
		query.Add("user_id", "1")
		url.RawQuery = query.Encode()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url.String(), nil)
		g.ServeHTTP(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("/douyin/user/: want status code: %v, but got %v", http.StatusOK, w.Result().StatusCode)
		}

		resp := parseUserResponse(w.Body, t)
		if resp.Response.StatusCode != 0 {
			t.Errorf("/douyin/user/: want response code: %v, but got %v", http.StatusOK, resp.Response.StatusCode)
		}
	})

	t.Run("user not exists", func(t *testing.T) {
		token, err := jwt.GenToken(1)
		if err != nil {
			t.Fatal(err.Error())
		}

		url := url.URL{Path: "/douyin/user/"}
		query := url.Query()
		query.Add("token", token)
		query.Add("user_id", "114")
		url.RawQuery = query.Encode()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url.String(), nil)
		g.ServeHTTP(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("/douyin/user/: want status code: %v, but got %v", http.StatusOK, w.Result().StatusCode)
		}

		resp := parseUserResponse(w.Body, t)
		if resp.Response.StatusCode != http.StatusBadRequest {
			t.Errorf("/douyin/user/: want response code: %v, but got %v", http.StatusBadRequest, resp.Response.StatusCode)
		}
	})

	t.Run("unauth", func(t *testing.T) {
		token := "invalid.token"

		url := url.URL{Path: "/douyin/user/"}
		query := url.Query()
		query.Add("token", token)
		query.Add("user_id", "1")
		url.RawQuery = query.Encode()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url.String(), nil)
		g.ServeHTTP(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("/douyin/user/: want status code: %v, but got %v", http.StatusOK, w.Result().StatusCode)
		}

		resp := parseUserResponse(w.Body, t)
		if resp.Response.StatusCode != http.StatusUnauthorized { //! not a big deal
			t.Errorf("/douyin/user/: want response code: %v, but got %v", http.StatusUnauthorized, resp.Response.StatusCode)
		}
	})
}
