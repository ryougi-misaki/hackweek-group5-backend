package middleware

import (
	"github.com/gin-gonic/gin"
	"hackathon/dao/mysql"
	"hackathon/models"
	"hackathon/response"
	"hackathon/util"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		//validate token format
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			response.Fail(ctx, nil, "权限不足")
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := util.ParseToken(tokenString)
		if err != nil || !token.Valid {
			response.Fail(ctx, nil, "权限不足")
			ctx.Abort()
			return
		}

		//验证通过后获取claims中的userid
		userId := claims.UserId
		DB := mysql.GetDB()
		var user models.User
		DB.First(&user, userId)

		//用户不存在
		if user.ID == 0 {
			response.Fail(ctx, nil, response.GetErrMsg(response.CodeUserNotExist))
			ctx.Abort()
			return
		}

		//用户存在,将user的信息写入上下文
		ctx.Set("user", user)

	}
}
