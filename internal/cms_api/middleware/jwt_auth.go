package middleware

import (
	"micro_servers/pkg/common/constant"
	"micro_servers/pkg/common/http"
	"micro_servers/pkg/common/log"
	"micro_servers/pkg/common/token_verify"
	"micro_servers/pkg/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), "")
		log.NewInfo("0", utils.GetSelfFuncName(), "userID: ", userID)
		c.Set("userID", userID)
		if !ok {
			log.NewError("", "GetUserIDFromToken false ", c.Request.Header.Get("token"))
			c.Abort()
			http.RespHttp200(c, constant.ErrParseToken, nil)
			return
		} else {
			log.NewInfo("0", utils.GetSelfFuncName(), "failed: ", errInfo)
		}
	}
}
