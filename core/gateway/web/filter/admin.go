package filter

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/project-template/cache/L2"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/bo"
	"github.com/project-template/errorcode"
	"net/http"
)

func AdminVerify(context *gin.Context) {
	if completeCoverage[context.Request.URL.Path] {
		return
	}
	resp := adminVerifyFunction(context)
	if resp.Code != errorcode.Success {
		// 没有权限，管理员被删除，冻结，电话，密码发生变更等，需要重新登录
		if resp.Code == errorcode.Unauthorized {
			context.Abort()
			context.JSON(http.StatusUnauthorized, resp.Reply(nil))
			return
		}
		// 当前账户名称发生变化，需要使用 long token 刷新
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

func adminVerifyFunction(ctx *gin.Context) *enp.Response {
	var ac = new(bo.AdminClaim)
	value, exists := ctx.Get(bo.AdminContextKey)
	if !exists {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("admin verify context not exists"))
	}
	ac, ok := value.(*bo.AdminClaim)
	if !ok {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("admin verify context not match"))
	}
	// 管理员信息正常检查
	admin, resp := L2.GetAdminById(ac.Id)
	if resp.Code != errorcode.Success {
		return resp
	}
	if admin == nil || admin.Id == 0 {
		return enp.Put(errorcode.GetAdminByIdNil)
	}
	if ac.UserName != admin.UserName {
		return enp.Put(errorcode.AdminTokenExpire)
	}
	if ac.Phone != admin.Phone {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("管理员电话改变"))
	}
	if ac.Password != admin.Password {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("管理员密码改变"))
	}
	if admin.IsDelete == true {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("管理员账号被删除"))
	}
	if admin.IsLock == true {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("管理员账号被锁定"))
	}
	// 单点登录检查
	if !L2.IsAdminSSo(ac.Id, ac.SSO) {
		return enp.Put(errorcode.Unauthorized, enp.AddIn(ac), enp.FormatMsg("管理员在其他地方登录"))
	}
	// ip 设置
	admin.RealIP = ctx.ClientIP()
	bytes, err := json.Marshal(admin)
	if err != nil {
		return enp.Put(errorcode.JsonMarshal, enp.AddIn(admin), enp.AddError(err))
	}
	ctx.Set(bo.HolderInformation, string(bytes))
	return enp.Put(errorcode.Success)
}
