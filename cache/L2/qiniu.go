package L2

import (
	"fmt"
	"github.com/project-template/common/config"
	"github.com/project-template/common/partner/qiniu"
	"time"
)

func GetQiniuToken() string {
	var key = fmt.Sprintf("tp_qiniu_token")
	val := config.Info().RedisClient.Get(key).Val()
	if len(val) != 0 {
		return val
	}
	token := qiniu.Token()
	config.Info().RedisClient.Set(key, token, 25*time.Minute)
	return token
}
