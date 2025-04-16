package filter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/project-template/cache/L2"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/bo"
	"github.com/project-template/errorcode"
	"net/http"
	"strconv"
	"strings"
)

func RoleVerify(context *gin.Context) {
	if completeCoverage[context.Request.URL.Path] {
		return
	}
	// 权限检查
	resp := roleVerifyFunction(context)
	if resp.Code != errorcode.Success {
		// 没有权限，管理员被删除，冻结，电话，密码发生变更等，需要重新登录
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

func roleVerifyFunction(ctx *gin.Context) *enp.Response {
	var ac = new(bo.AdminClaim)
	value, exists := ctx.Get(bo.AdminContextKey)
	if !exists {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("role verify context not exists"))
	}
	ac, ok := value.(*bo.AdminClaim)
	if !ok {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("role verify context not match"))
	}
	roleId, resp := L2.GetAdminRole(ac.Id)
	if resp.Code != errorcode.Success {
		return resp
	}
	if roleId == 0 {
		return enp.Put(errorcode.GetAdminRoleNil, enp.AddIn(ac.Id))
	}
	role, resp := L2.GetRoleById(roleId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if role == nil || role.Id == 0 {
		return enp.Put(errorcode.GetRoleByIdNil, enp.AddIn(roleId))
	}
	var menuStr, pageStr, buttonStr = role.MenuStr, role.PageStr, role.ButtonStr
	// 管理员权限是否改变
	if isMenuSame, isPageSame, isButtonSame := isRoleSame(menuStr, ac.MenuStr), isRoleSame(pageStr, ac.PageStr), isRoleSame(buttonStr, ac.ButtonStr); !isMenuSame || !isPageSame || !isButtonSame {
		return enp.Put(errorcode.AdminTokenExpire, enp.FormatMsg("role verify context not match"))
	}

	// 该接口是否有权限访问
	if resp = isRolePass(ctx.Request.Method, ctx.Request.URL.Path, menuStr, pageStr, buttonStr); resp.Code != errorcode.Success {
		return resp
	}

	return enp.Put(errorcode.Success)
}

func isRoleSame(param, source string) bool {
	prs := strings.Split(param, ",")
	srs := strings.Split(source, ",")
	if len(prs) == 0 && len(srs) == 0 {
		// 例如 button str 修改前后都是空，那么直接是true
		return true
	}
	if len(prs) == 0 || len(srs) == 0 {
		// 通过上面判断后，那么只会存在，一个是空，一个有值，必然不同
		return false
	}
	prm := make(map[string]bool)
	for _, pr := range prs {
		prm[pr] = true
	}
	srm := make(map[string]bool)
	for _, sr := range srs {
		srm[sr] = true
	}
	for k := range prm {
		if !srm[k] {
			return false
		}
	}
	for k := range srm {
		if !prm[k] {
			return false
		}
	}
	return true
}

func isRolePass(method, path, menuStr, pageStr, buttonStr string) *enp.Response {
	limit := config.Info().ApiLimit[fmt.Sprintf("%s$####$%s", method, path)]
	if limit == nil {
		// 不受权限控制
		return enp.Put(errorcode.Success)
	}
	var menuId, pageId, buttonId = limit.MenuId, limit.PageId, limit.ButtonId
	if menuId == 0 && pageId == 0 && buttonId == 0 {
		// 全部为0，不受权限控制
		return enp.Put(errorcode.Success)
	}

	isMenuPass, resp := isMatch(menuId, menuStr)
	if resp.Code != errorcode.Success {
		return resp
	}
	isPagePass, resp := isMatch(pageId, pageStr)
	if resp.Code != errorcode.Success {
		return resp
	}
	isButtonPass, resp := isMatch(buttonId, buttonStr)
	if resp.Code != errorcode.Success {
		return resp
	}
	if !isMenuPass || !isPagePass || !isButtonPass {
		// 只要一项没过，全部不过
		return enp.Put(errorcode.RolePass)
	}
	return enp.Put(errorcode.Success)
}

func isMatch(roleId int64, rst string) (bool, *enp.Response) {
	if roleId == 0 || len(rst) == 0 {
		// 当menuId==0时，验证标识直接为true，不为0时就需要验证具体权限值，page，button同理
		return true, enp.Put(errorcode.Success)
	}
	bs := strings.Split(rst, ",")
	for _, b := range bs {
		ri, err := strconv.ParseInt(b, 10, 64)
		if err != nil {
			return false, enp.Put(errorcode.StrconvParseInt, enp.AddError(err))
		}
		if ri == roleId {
			return true, enp.Put(errorcode.Success)
		}
	}
	return false, enp.Put(errorcode.RolePass)
}
