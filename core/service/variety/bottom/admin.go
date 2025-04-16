package bottom

import (
	"database/sql"
	"errors"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/errorcode"
)

func GetAdminById(id int64) (*do.Admin, *enp.Response) {
	a := new(do.Admin)
	err := config.Info().MysqlClient.QueryRow("SELECT "+
		"`id`,`username`,`phone`,`password`,`is_delete`,`is_lock`,`create_time`,`update_time` "+
		"FROM `trip_portal`.`admin` WHERE `id`=?", id).Scan(
		&a.Id, &a.UserName, &a.Phone, &a.Password, &a.IsDelete, &a.IsLock, &a.CreateTime, &a.UpdateTime)
	if err != nil {
		return nil, enp.Put(errorcode.MysqlScanErr)
	}
	return a, enp.Put(errorcode.Success)
}

func GetAdminByPhone(phone string) (int64, *enp.Response) {
	var adminId int64
	err := config.Info().MysqlClient.QueryRow("SELECT `id` FROM `trip_portal`.`admin` WHERE `phone`=? AND `is_delete`=false", phone).Scan(&adminId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, enp.Put(errorcode.Success, enp.AddIn(phone))
	}
	if err != nil {
		return 0, enp.Put(errorcode.MysqlScanErr, enp.AddError(err))
	}
	return adminId, enp.Put(errorcode.Success)
}
