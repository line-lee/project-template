package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/project-template/common/client/variety"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	mw "github.com/project-template/common/middleware"
	"github.com/project-template/common/quit"
	"github.com/project-template/common/tools"
	"github.com/project-template/core/gateway/web/filter"
	"github.com/project-template/errorcode"
	"net/http"
)

func main() {
	// 初始化配置信息和中间件链接
	config.Init().Open(mw.OpenMysqlConnect, mw.OpenRedisConnect)
	// 注册服务关闭时，执行的操作事务
	quit.GetQuitEvent().RegisterFunc(mw.CloseMysqlConnect, mw.CloseRedisConnect)
	// 开启gin路由
	start()
	// 优雅退出
	quit.WaitSignal()
}

func start() {
	tools.SecureGo(func(args ...interface{}) {
		engine := gin.Default()
		// 健康检查
		engine.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
		// web路由
		webGroup := engine.Group("/web/api/v1")
		{
			webGroup.Use(filter.TokenVerify, filter.AdminVerify, filter.RoleVerify)
			// API-Document
			variety.UtilsRouter(webGroup)
			variety.AdminRouter(webGroup)
			variety.RoleRouter(webGroup)
		}
		// open路由 .......................

		if err := engine.Run(fmt.Sprintf(":%d", config.Info().Core.Gateways[config.WebGateway].Http)); err != nil {
			enp.Put(errorcode.GinRunErr)
		}
	})
}
