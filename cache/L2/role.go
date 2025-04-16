package L2

import (
	"encoding/json"
	"fmt"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	mw "github.com/project-template/common/middleware"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/core/service/variety/bottom"
	"github.com/project-template/errorcode"
	"strconv"
	"time"
)

func GetRoleById(id int64) (*do.Role, *enp.Response) {
	var key = fmt.Sprintf("tp_role_%d", id)
	val := config.Info().RedisClient.Get(key).Val()
	if len(val) != 0 {
		if val == mw.RedisDefaultValue {
			return nil, enp.Put(errorcode.GetRoleByIdNil, enp.AddIn(id))
		}
		role := new(do.Role)
		err := json.Unmarshal([]byte(val), role)
		if err != nil {
			return nil, enp.Put(errorcode.GetRoleByIdNil, enp.AddIn(val), enp.AddError(err))
		}
		return role, enp.Put(errorcode.Success)
	}
	role, resp := bottom.GetRoleById(id)
	if resp.Code != errorcode.Success {
		return nil, resp
	}
	if role == nil || role.Id == 0 {
		// 数据库查不到信息，就不要一直放下去查了，这样会击穿数据库
		// 同样的，在新增的时候要做删除缓存的操作，避免偶然碰撞事件发生
		config.Info().RedisClient.Set(key, mw.RedisDefaultValue, 10*time.Hour)
		return nil, enp.Put(errorcode.Success, enp.AddIn(id))
	}
	config.Info().RedisClient.Set(key, role, 10*time.Minute)
	return role, enp.Put(errorcode.Success)
}

func DelRole(id int64) *enp.Response {
	var key = fmt.Sprintf("tp_role_%d", id)
	config.Info().RedisClient.Del(key)
	return enp.Put(errorcode.Success)
}

func GetAdminRole(adminId int64) (int64, *enp.Response) {
	var key = fmt.Sprintf("tp_adminrole_%d", adminId)
	val := config.Info().RedisClient.Get(key).Val()
	if len(val) > 0 {
		roleId, _ := strconv.ParseInt(val, 10, 64)
		if roleId != 0 {
			return roleId, enp.Put(errorcode.Success)
		}
	}
	adminRole, resp := bottom.GetAdminRole(adminId)
	if resp.Code != errorcode.Success {
		return 0, resp
	}
	if adminRole.RoleId == 0 {
		// 数据库查不到信息，就不要一直放下去查了，这样会击穿数据库
		// 同样的，在新增的时候要做删除缓存的操作，避免偶然碰撞事件发生
		config.Info().RedisClient.Set(key, mw.RedisDefaultValue, 10*time.Hour)
		return 0, enp.Put(errorcode.Success, enp.AddIn(adminId))
	}
	config.Info().RedisClient.Set(key, adminRole.RoleId, 10*time.Minute)
	return adminRole.RoleId, enp.Put(errorcode.Success)
}

func DelAdminRole(adminId int64) {
	var key = fmt.Sprintf("tp_adminrole_%d", adminId)
	config.Info().RedisClient.Del(key)
}
