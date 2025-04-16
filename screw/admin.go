package screw

import (
	"bytes"
	"database/sql"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/common/tools"
	"github.com/project-template/errorcode"
)

func QueryRolesByAdminIds(adminIds string) (map[int64]*do.Role, *enp.Response) {
	db := config.Info().MysqlClient
	// 获取管理员id和权限id的映射列表
	adminRolesRows, err := db.Query("SELECT `admin_id`,`role_id` FROM `trip_portal`.`admin_roles` WHERE `admin_id` IN (" + adminIds + ")")
	if err != nil {
		return nil, enp.Put(errorcode.MysqlQueryErr, enp.AddIn(adminIds), enp.AddError(err))
	}
	defer func(adminRolesRows *sql.Rows) {
		err = adminRolesRows.Close()
		if err != nil {
			enp.Put(errorcode.MysqlRowsCloseErr)
		}
	}(adminRolesRows)
	var roleIds bytes.Buffer
	var roleFilter = make(map[int64]bool)
	var arm = make(map[int64]int64)
	for adminRolesRows.Next() {
		ar := new(do.AdminRole)
		err = adminRolesRows.Scan(&ar.AdminId, &ar.RoleId)
		if err != nil {
			return nil, enp.Put(errorcode.MysqlScanErr, enp.AddIn(adminIds), enp.AddError(err))
		}
		if !roleFilter[ar.RoleId] {
			roleIds = tools.ConcatWith(ar.RoleId, roleIds, tools.ConcatWithComma)
			roleFilter[ar.RoleId] = true
		}
		arm[ar.AdminId] = ar.RoleId
	}
	roleRows, err := db.Query("SELECT `id`,`name` FROM `trip_portal`.`roles` WHERE `id` IN (" + roleIds.String() + ")")
	if err != nil {
		return nil, enp.Put(errorcode.MysqlQueryErr, enp.AddIn(roleIds.String()), enp.AddError(err))
	}
	defer func(roleRows *sql.Rows) {
		err := roleRows.Close()
		if err != nil {
			enp.Put(errorcode.MysqlRowsCloseErr)
		}
	}(roleRows)
	var rim = make(map[int64]*do.Role)
	for roleRows.Next() {
		role := new(do.Role)
		err = roleRows.Scan(&role.Id, &role.Name)
		if err != nil {
			return nil, enp.Put(errorcode.MysqlScanErr, enp.AddIn(roleIds.String()), enp.AddError(err))
		}
		rim[role.Id] = role
	}
	rm := make(map[int64]*do.Role)
	for k, v := range arm {
		rm[k] = rim[v]
	}
	return rm, enp.Put(errorcode.Success)
}
