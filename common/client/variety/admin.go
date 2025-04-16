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

func admin(rout adminRout, param interface{}) *proto.Request {
	request := &proto.Request{
		Service: config.VarietyService,
		Module:  adminModule,
		Api:     int64(rout),
	}
	b, _ := json.Marshal(param)
	request.Data = b
	return request
}

type adminRout int

const (
	// *****************************************************************
	// 特别注意：如果调整了顺序，需要同时更新服务的发现和注册
	// *****************************************************************
	adminLoginRout        adminRout = 1
	adminTokenRefreshRout adminRout = 2
)

var adminFunctions = map[adminRout]func([]byte) *enp.Response{
	adminLoginRout:        impl.AdminLogin,
	adminTokenRefreshRout: impl.AdminTokenRefresh,
}

func AdminRouter(group *gin.RouterGroup) {
	// 管理员登录
	group.POST("/login", adminLogin)
	// 管理员token刷新
	group.POST("/token/refresh", adminTokenRefresh)
}

func adminLogin(context *gin.Context) {
	param := new(vo.AdminLoginParam)
	if err := tools.GinParamBind(context, param); err != nil {
		context.JSON(http.StatusBadRequest, enp.Put(errorcode.GinShouldBindErr).Reply(nil))
		return
	}
	// 客户端ip设置
	param.RealIP = context.ClientIP()
	context.JSON(http.StatusOK, proto.GRPC(context, admin(adminLoginRout, param)).Reply(new(vo.AdminLoginResponse)))
}

func adminTokenRefresh(context *gin.Context) {
	param := new(vo.AdminTokenRefreshParam)
	if err := tools.GinParamBind(context, param); err != nil {
		context.JSON(http.StatusBadRequest, enp.Put(errorcode.GinShouldBindErr).Reply(nil))
		return
	}
	responses := new(vo.AdminTokenRefreshResponse)
	reply := proto.GRPC(context, admin(adminTokenRefreshRout, param)).Reply(responses)
	if reply.Code == errorcode.Unauthorized {
		context.JSON(http.StatusUnauthorized, reply)
		return
	}
	context.JSON(http.StatusOK, reply)
}
