package qiniu

import (
	"github.com/project-template/common/config"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func Token() string {
	putPolicy := storage.PutPolicy{
		Scope: config.Info().QiNiuConfig.Bucket,
	}
	putPolicy.Expires = 30 * 60 //30分钟
	mac := qbox.NewMac(config.Info().QiNiuConfig.AccessKey, config.Info().QiNiuConfig.SecretKey)
	token := putPolicy.UploadToken(mac)

	return token
}
