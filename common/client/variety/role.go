package variety

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/project-template/common/client/proto"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/vo"
	"github.com/project-template/common/tools"
	"github.com/project-template/core/service/variety/impl"
	"github.com/project-template/errorcode"
	"net/http"
)

func role(rout roleRout, param interface{}) *proto.Request {
	request := &proto.Request{
		Service: config.VarietyService,
		Module:  roleModule,
		Api:     int64(rout),
	}
	b, _ := json.Marshal(param)
	request.Data = b
	return request
}

type roleRout int

const (
	// *****************************************************************
	// 特别注意：如果调整了顺序，需要同时更新服务的发现和注册
	// *****************************************************************
	queryAllAuthsRout roleRout = 1
	addRoleRout       roleRout = 2
)

var roleFunctions = map[roleRout]func([]byte) *enp.Response{
	queryAllAuthsRout: impl.QueryAllAuths,
	addRoleRout:       impl.AddRole,
}

func RoleRouter(group *gin.RouterGroup) {
	// 查询系统全部权限集合
	group.GET("/auths", queryAllAuths)
	// 新增角色
	group.POST("/role", addRole)

}

func queryAllAuths(context *gin.Context) {
	param := new(vo.QueryAllAuthsParam)
	if err := tools.GinParamBind(context, param); err != nil {
		context.JSON(http.StatusBadRequest, enp.Put(errorcode.GinShouldBindErr).Reply(nil))
		return
	}
	context.JSON(http.StatusOK, proto.GRPC(context, role(queryAllAuthsRout, param)).Reply(new(vo.QueryAllAuthsResponse)))
}

func addRole(context *gin.Context) {
	param := new(vo.AddRoleParam)
	if err := tools.GinParamBind(context, param); err != nil {
		context.JSON(http.StatusBadRequest, enp.Put(errorcode.GinShouldBindErr).Reply(nil))
		return
	}
	context.JSON(http.StatusOK, proto.GRPC(context, role(addRoleRout, param)).Reply(nil))
}
