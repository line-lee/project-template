package L2

import (
	"encoding/json"
	"fmt"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	mw "github.com/project-template/common/middleware"
	"github.com/project-template/common/models/variety/bo"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/core/service/variety/bottom"
	"github.com/project-template/errorcode"
	"strconv"
	"time"
)

func GetAdminById(id int64) (*do.Admin, *enp.Response) {
	var key = fmt.Sprintf("tp_admin_%d", id)
	val := config.Info().RedisClient.Get(key).Val()
	if len(val) != 0 {
		if val == mw.RedisDefaultValue {
			return nil, enp.Put(errorcode.GetAdminByIdNil, enp.AddIn(id))
		}
		admin := new(do.Admin)
		err := json.Unmarshal([]byte(val), admin)
		if err != nil {
			return nil, enp.Put(errorcode.GetAdminByIdNil, enp.AddIn(val), enp.AddError(err))
		}
		return admin, enp.Put(errorcode.Success)
	}
	admin, resp := bottom.GetAdminById(id)
	if resp.Code != errorcode.Success {
		return nil, resp
	}
	if admin == nil || admin.Id == 0 {
		// 数据库查不到信息，就不要一直放下去查了，这样会击穿数据库
		// 同样的，在新增的时候要做删除缓存的操作，避免偶然碰撞事件发生
		config.Info().RedisClient.Set(key, mw.RedisDefaultValue, 10*time.Hour)
		return nil, enp.Put(errorcode.Success)
	}
	config.Info().RedisClient.Set(key, admin, 10*time.Minute)
	return admin, enp.Put(errorcode.Success)
}

func DelAdminById(id int64) {
	var key = fmt.Sprintf("tp_admin_%d", id)
	config.Info().RedisClient.Del(key)
}

func GetAdminByPhone(phone string) (int64, *enp.Response) {
	var key = fmt.Sprintf("tp_admin_phone_%s", phone)
	adminId := config.Info().RedisClient.Get(key).Val()
	if len(adminId) != 0 {
		if adminId == mw.RedisDefaultValue {
			return 0, enp.Put(errorcode.Success, enp.AddIn(phone))
		}
		id, err := strconv.ParseInt(adminId, 10, 64)
		if err != nil {
			return 0, enp.Put(errorcode.StrconvParseInt, enp.AddIn(adminId))
		}
		return id, enp.Put(errorcode.Success)
	}
	id, resp := bottom.GetAdminByPhone(phone)
	if resp.Code != errorcode.Success {
		return 0, resp
	}
	if id == 0 {
		// 数据库查不到信息，就不要一直放下去查了，这样会击穿数据库
		// 同样的，在新增的时候要做删除缓存的操作，避免偶然碰撞事件发生
		config.Info().RedisClient.Set(key, mw.RedisDefaultValue, 10*time.Hour)
		return 0, enp.Put(errorcode.Success, enp.AddIn(id))
	}
	config.Info().RedisClient.Set(key, id, 10*time.Minute)
	return id, enp.Put(errorcode.Success)
}

func DelAdminByPhone(phone string) {
	var key = fmt.Sprintf("tp_admin_phone_%s", phone)
	config.Info().RedisClient.Del(key)
}

func IsAdminSSo(id int64, sso string) bool {
	var key = fmt.Sprintf("tp_admin_sso_%d", id)
	str := config.Info().RedisClient.Get(key).Val()
	return str == sso
}

func SetAdminSSo(id int64, sso string) {
	var key = fmt.Sprintf("tp_admin_sso_%d", id)
	config.Info().RedisClient.Set(key, sso, bo.AdminTokenLongExpire)
}
