package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/project-template/common/client"
	"github.com/project-template/common/client/proto"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	mw "github.com/project-template/common/middleware"
	"github.com/project-template/common/quit"
	"github.com/project-template/common/tools"
	"github.com/project-template/core/service/variety/listener"
	"github.com/project-template/errorcode"
	"net/http"
)

func main() {
	// 初始化配置信息和中间件链接
	config.Init().Open(mw.OpenMysqlConnect, mw.OpenRedisConnect, mw.OpenKafka)
	// 注册服务关闭时，执行的操作事务
	quit.GetQuitEvent().RegisterFunc(mw.CloseMysqlConnect, mw.CloseRedisConnect, mw.CloseKafka)
	// kafka监听
	listener.Start()
	// GRPC注册启动
	proto.GRPCStart(config.Info().Core.Services[config.VarietyService].Grpc, &service{})
	// 开启gin
	ginStart()
	// 优雅退出
	quit.WaitSignal()
}

type service struct{}

func (s *service) REQ(_ context.Context, in *proto.Request) (resp *proto.Response, err error) {
	defer func() {
		if v := recover(); v != nil {
			bytes, _ := json.Marshal(enp.Put(errorcode.Recover, enp.AddOut(v)))
			resp = &proto.Response{Data: bytes}
		}
	}()
	return client.Do(in.Service, int(in.Module), int(in.Api), in.Data)
}

func ginStart() {
	tools.SecureGo(func(args ...interface{}) {
		engine := gin.Default()
		engine.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
		err := engine.Run(fmt.Sprintf(":%d", config.Info().Core.Services[config.VarietyService].Http))
		if err != nil {
			enp.Put(errorcode.GinRunErr)
		}
	})
}
