package variety

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/project-template/cache/L2"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/vo"
	"github.com/project-template/common/tools"
	"github.com/project-template/errorcode"
	"net/http"
)

func UtilsRouter(group *gin.RouterGroup) {
	// 百度地图三级行政区域数据
	group.GET("/baidu/location", baiduLocation)
	// 密码生成器（16位）
	group.POST("/password/generate", passwordGenerate)
	// 获取七牛云token
	group.GET("/qiniu/token", qiniuToken)
}

func baiduLocation(context *gin.Context) {
	provinces := make([]*vo.Province, 0)
	err := json.Unmarshal([]byte(config.Info().BaiduLocation), &provinces)
	if err != nil {
		context.JSON(http.StatusOK, enp.Put(errorcode.JsonUnmarshal, enp.AddError(err)).Reply(nil))
	}
	resp := vo.BaiduLocationResponse{Provinces: provinces}
	context.JSON(http.StatusOK, enp.Put(errorcode.Success, enp.AddData(resp)).Reply(new(vo.AdminLoginResponse)))
}

func passwordGenerate(context *gin.Context) {
	resp := vo.PasswordGenerateResponse{Password: tools.RandomString(16, tools.AnyMod)}
	context.JSON(http.StatusOK, enp.Put(errorcode.Success, enp.AddData(resp)).Reply(new(vo.PasswordGenerateResponse)))
}

func qiniuToken(context *gin.Context) {
	resp := vo.GetQiNiuTokenResponse{Token: L2.GetQiniuToken()}
	context.JSON(http.StatusOK, enp.Put(errorcode.Success, enp.AddData(resp)).Reply(new(vo.GetQiNiuTokenResponse)))
}
