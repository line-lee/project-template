package bottom

import (
	"database/sql"
	"errors"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/errorcode"
)

func GetRoleById(id int64) (*do.Role, *enp.Response) {
	role := new(do.Role)
	err := config.Info().MysqlClient.QueryRow("select `id`,`name`,`description`,`menu_str`,`page_str`,`button_str`,"+
		"`is_main`,`is_deleted`,`updated`,`created`,`version` from `trip_portal`.`roles` where `id` = ?", id).Scan(
		&role.Id, &role.Name, &role.Description, &role.MenuStr, &role.PageStr, &role.ButtonStr, &role.IsMain,
		&role.IsDeleted, &role.Updated, &role.Created, &role.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, enp.Put(errorcode.Success, enp.AddIn(id))
	}
	if err != nil {
		return nil, enp.Put(errorcode.MysqlScanErr, enp.AddIn(id), enp.AddError(err))
	}
	return role, enp.Put(errorcode.Success)
}

func GetAdminRole(adminId int64) (*do.AdminRole, *enp.Response) {
	ar := new(do.AdminRole)
	err := config.Info().MysqlClient.QueryRow("select `admin_id`,`role_id` from `trip_portal`.`admin_roles` where `admin_id` = ?", adminId).Scan(&ar.AdminId, &ar.RoleId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, enp.Put(errorcode.Success, enp.AddIn(adminId))
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, enp.Put(errorcode.MysqlScanErr, enp.AddIn(adminId), enp.AddError(err))
	}
	return ar, enp.Put(errorcode.Success)
}
