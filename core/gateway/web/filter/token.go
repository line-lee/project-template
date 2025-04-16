package filter

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/bo"
	"github.com/project-template/common/tools"
	"github.com/project-template/errorcode"
	"net/http"
	"time"
)

func TokenVerify(context *gin.Context) {
	if completeCoverage[context.Request.URL.Path] {
		return
	}
	resp := tokenVerifyFunction(context)
	if resp.Code != errorcode.Success {
		// 没有权限，或者长短token都失效的时候直接不通过，需要重新登录
		if resp.Code == errorcode.Unauthorized {
			context.Abort()
			context.JSON(http.StatusUnauthorized, resp.Reply(nil))
			return
		}
		// 当前 short token 失效，需要使用 long token 刷新
		if resp.Code == errorcode.AdminTokenExpire {
			context.Abort()
			context.JSON(http.StatusOK, resp.Reply(nil))
			return
		}
		// 当遇到一些系统连接错误，比如查询管理员查不到（可能是grpc连接问题），我们直接返回错误，但是不会返回登录，也不会刷新token
		context.Abort()
		context.JSON(http.StatusOK, resp.Reply(nil))
		return
	}
}

func tokenVerifyFunction(ctx *gin.Context) *enp.Response {
	var thisTime = time.Now()
	tokenStr := ctx.Request.Header.Get(bo.AdminTokenKey)
	if len(tokenStr) == 0 {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("token 为空"))
	}
	// aes解密
	decrypt, err := tools.AESDecrypt(tokenStr)
	if err != nil {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("aes 解密失败"))
	}
	ac := new(bo.AdminClaim)
	err = json.Unmarshal(decrypt, ac)
	if err != nil {
		return enp.Put(errorcode.JsonMarshal, enp.AddIn(string(decrypt)), enp.AddError(err))
	}
	if ac.Id == 0 {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("token id 0"))
	}
	// 过期检查
	if !tokenExpire[ctx.Request.URL.Path] && thisTime.Unix() > ac.ExpireTime {
		return enp.Put(errorcode.AdminTokenExpire, enp.AddIn(thisTime.Unix(), ac))
	}
	// 类型检查
	if ac.Type != bo.TokenShort {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("token not short"))
	}
	ctx.Set(bo.AdminContextKey, ac)
	return enp.Put(errorcode.Success)
}
