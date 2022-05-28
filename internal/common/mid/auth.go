package mid

import (
	"net/http"

	"github.com/fitenne/youthcampus-dousheng/internal/common/jwt"
	"github.com/fitenne/youthcampus-dousheng/internal/controller"
	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "userID"
)

// JWTAuthMiddleware 在需要判断用户登录状态的地方使用此中间件，若判断登录状态不合法，则直接返回响应，且不会执行后续的处理函数
// example: api.GET("/xxx/xx/", mid.JWTAuthMiddleware(),xxxxHandler)
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//  Token放在query中
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}
		if token == "" {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: http.StatusUnauthorized,
				StatusMsg:  "invalid token",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: http.StatusUnauthorized,
				StatusMsg:  "invalid token",
			})
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(UserIDKey, mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("userID")来获取当前请求的用户信息
	}

}
