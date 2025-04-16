package impl

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/project-template/cache/L2"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/bo"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/common/models/variety/vo"
	"github.com/project-template/common/tools"
	"github.com/project-template/errorcode"
	"github.com/project-template/screw"
	"github.com/project-template/screw/log"
	"time"
)

func AdminLogin(data []byte) *enp.Response {
	param := new(vo.AdminLoginParam)
	err := json.Unmarshal(data, param)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddIn(string(data)), enp.AddError(err))
	}
	if len(param.Phone) == 0 {
		return enp.Put(errorcode.InvalidParam, enp.AddIn(param), enp.FormatMsg("Phone", param.Phone))
	}
	if len(param.Password) == 0 {
		return enp.Put(errorcode.InvalidParam, enp.AddIn(param), enp.FormatMsg("Password", param.Password))
	}
	var logWriteParam = &log.AdminLoginLogWriterParam{Param: param}
	// 根据手机号密码查询账户是否存在
	adminId, resp := L2.GetAdminByPhone(param.Phone)
	if resp.Code != errorcode.Success {
		return resp
	}
	if adminId == 0 {
		return enp.Put(errorcode.AdminPassword, enp.AddIn(param.Phone))
	}
	// 管理员详细信息
	admin, resp := L2.GetAdminById(adminId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if admin == nil || admin.Id == 0 {
		return enp.Put(errorcode.AdminPassword, enp.AddIn(adminId))
	}
	if admin.IsDelete {
		return enp.Put(errorcode.AdminDelete, enp.AddIn(admin))
	}
	if admin.IsLock {
		return enp.Put(errorcode.AdminLock, enp.AddIn(admin))
	}
	if admin.Password != param.Password {
		return enp.Put(errorcode.AdminPassword, enp.AddIn(admin))
	}
	logWriteParam.Admin = admin
	// 获取管理员权限
	roleId, resp := L2.GetAdminRole(adminId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if roleId == 0 {
		return enp.Put(errorcode.GetAdminRoleNil, enp.AddIn(adminId))
	}
	role, resp := L2.GetRoleById(roleId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if role == nil || role.Id == 0 {
		return enp.Put(errorcode.GetRoleByIdNil, enp.AddIn(roleId))
	}
	// 生成token
	thisTime := time.Now()
	sso := uuid.New().String()
	claim := bo.AdminClaim{
		Id:        adminId,
		UserName:  admin.UserName,
		Phone:     admin.Phone,
		Password:  admin.Password,
		SSO:       sso,
		RoleId:    roleId,
		MenuStr:   role.MenuStr,
		PageStr:   role.PageStr,
		ButtonStr: role.ButtonStr,
	}
	shortClaim := claim
	shortClaim.Type = bo.TokenShort
	shortExpire := thisTime.Add(bo.AdminTokenShortExpire)
	shortClaim.ExpireTime = shortExpire.Unix()
	shortBytes, err := json.Marshal(shortClaim)
	if err != nil {
		return enp.Put(errorcode.JsonMarshal, enp.AddIn(shortClaim))
	}
	shortToken, err := tools.AESEncrypt(shortBytes)
	if err != nil {
		return enp.Put(errorcode.AESEncrypt, enp.AddIn(string(shortBytes), enp.AddError(err)))
	}
	longClaim := claim
	longClaim.Type = bo.TokenLong
	longExpire := thisTime.Add(bo.AdminTokenLongExpire)
	longClaim.ExpireTime = longExpire.Unix()
	longBytes, err := json.Marshal(longClaim)
	if err != nil {
		return enp.Put(errorcode.JsonMarshal, enp.AddIn(longClaim))
	}
	longToken, err := tools.AESEncrypt(longBytes)
	if err != nil {
		return enp.Put(errorcode.AESEncrypt, enp.AddIn(string(longBytes), enp.AddError(err)))
	}
	response := vo.AdminLoginResponse{
		Name:         admin.UserName,
		Phone:        admin.Phone,
		Menu:         role.MenuStr,
		Page:         role.PageStr,
		Button:       role.ButtonStr,
		Token:        shortToken,
		RefreshToken: longToken,
	}
	L2.SetAdminSSo(adminId, sso)
	log.Add(logWriteParam)
	return enp.Put(errorcode.Success, enp.AddData(response))
}

func AdminTokenRefresh(data []byte) *enp.Response {
	param := new(vo.AdminTokenRefreshParam)
	err := json.Unmarshal(data, param)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddIn(string(data)), enp.AddError(err))
	}
	if len(param.RefreshToken) == 0 {
		return enp.Put(errorcode.InvalidParam, enp.AddIn(param), enp.FormatMsg("RefreshToken", param.RefreshToken))
	}
	tb, err := tools.AESDecrypt(param.RefreshToken)
	if err != nil {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("aes 解密错误"))
	}
	ac := new(bo.AdminClaim)
	err = json.Unmarshal(tb, ac)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddError(err))
	}
	if ac.Id == 0 {
		return enp.Put(errorcode.AdminTokenIdNil, enp.AddError(err))
	}
	// sso校验
	if !L2.IsAdminSSo(ac.Id, ac.SSO) {
		return enp.Put(errorcode.Unauthorized, enp.AddIn(ac), enp.FormatMsg("sso 不匹配"))
	}
	// refresh token 过期
	thisTime := time.Now()
	if ac.ExpireTime < thisTime.Unix() {
		return enp.Put(errorcode.Unauthorized, enp.AddIn(ac.ExpireTime, thisTime.Unix()), enp.FormatMsg("过期"))
	}
	// 管理员详细信息
	a, resp := L2.GetAdminById(ac.Id)
	if resp.Code != errorcode.Success {
		return resp
	}
	if a == nil || a.Id == 0 {
		return enp.Put(errorcode.GetAdminByIdNil)
	}
	if a.IsDelete {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("账号被删除"))
	}
	if a.IsLock {
		return enp.Put(errorcode.Unauthorized, enp.FormatMsg("账号被锁定"))
	}
	// 获取管理员权限
	roleId, resp := L2.GetAdminRole(a.Id)
	if resp.Code != errorcode.Success {
		return resp
	}
	if roleId == 0 {
		return enp.Put(errorcode.GetAdminRoleNil, enp.AddIn(a.Id))
	}
	role, resp := L2.GetRoleById(roleId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if role == nil || role.Id == 0 {
		return enp.Put(errorcode.GetRoleByIdNil, enp.AddIn(roleId))
	}
	// 生成token
	sso := uuid.New().String()
	claim := bo.AdminClaim{
		Id:        ac.Id,
		UserName:  a.UserName,
		Phone:     a.Phone,
		Password:  a.Password,
		SSO:       sso,
		RoleId:    roleId,
		MenuStr:   role.MenuStr,
		PageStr:   role.PageStr,
		ButtonStr: role.ButtonStr,
	}
	shortClaim := claim
	shortClaim.Type = bo.TokenShort
	shortExpire := thisTime.Add(bo.AdminTokenShortExpire)
	shortClaim.ExpireTime = shortExpire.Unix()
	shortBytes, err := json.Marshal(shortClaim)
	if err != nil {
		return enp.Put(errorcode.JsonMarshal, enp.AddIn(shortClaim))
	}
	shortToken, err := tools.AESEncrypt(shortBytes)
	if err != nil {
		return enp.Put(errorcode.AESEncrypt, enp.AddIn(string(shortBytes), enp.AddError(err)))
	}
	longClaim := claim
	longClaim.Type = bo.TokenLong
	longExpire := thisTime.Add(bo.AdminTokenLongExpire)
	longClaim.ExpireTime = longExpire.Unix()
	longBytes, err := json.Marshal(longClaim)
	if err != nil {
		return enp.Put(errorcode.JsonMarshal, enp.AddIn(longClaim))
	}
	longToken, err := tools.AESEncrypt(longBytes)
	if err != nil {
		return enp.Put(errorcode.AESEncrypt, enp.AddIn(string(longBytes), enp.AddError(err)))
	}
	response := vo.AdminTokenRefreshResponse{
		Name:         a.UserName,
		Phone:        a.Phone,
		Menu:         role.MenuStr,
		Page:         role.PageStr,
		Button:       role.ButtonStr,
		Token:        shortToken,
		RefreshToken: longToken,
	}
	L2.SetAdminSSo(a.Id, sso)
	return enp.Put(errorcode.Success, enp.AddData(response))
}

func QueryAdmins(data []byte) *enp.Response {
	param := new(vo.QueryAdminsParam)
	err := json.Unmarshal(data, param)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddError(err))
	}
	if param.Page == 0 {
		param.Page = 1
	}
	if param.Limit == 0 || param.Limit > 60 {
		param.Limit = 15
	}
	var countSql, querySql, commonSql bytes.Buffer
	var sqlParam = make([]any, 0)
	countSql.WriteString("SELECT COUNT(`id`) FROM `trip_portal`.`admin` WHERE `is_delete`=FALSE ")
	querySql.WriteString("SELECT `id`,`username`,`phone`,`create_time` FROM `trip_portal`.`admin` WHERE `is_delete`=FALSE ")
	if len(param.Phone) != 0 {
		commonSql.WriteString("AND `phone`=? ")
		sqlParam = append(sqlParam, param.Phone)
	}
	if len(param.Name) != 0 {
		commonSql.WriteString("AND `username` like ? ")
		sqlParam = append(sqlParam, "%"+param.Name+"%")
	}
	fmt.Println(commonSql.String())
	countSql.WriteString(commonSql.String())
	var total int64
	err = config.Info().MysqlClient.QueryRow(countSql.String(), sqlParam...).Scan(&total)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddIn(countSql.String(), sqlParam), enp.AddError(err))
	}
	response := new(vo.QueryAdminsResponse)
	response.Total = total
	if total == 0 {
		return enp.Put(errorcode.Success, enp.AddData(response))
	}
	querySql.WriteString(commonSql.String())
	querySql.WriteString("ORDER BY id DESC LIMIT ? OFFSET ? ")
	sqlParam = append(sqlParam, []interface{}{param.Limit, (param.Page - 1) * param.Limit}...)

	rows, err := config.Info().MysqlClient.Query(querySql.String(), sqlParam...)
	if err != nil {
		return enp.Put(errorcode.MysqlQueryErr, enp.AddIn(querySql.String(), sqlParam), enp.AddError(err))
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			enp.Put(errorcode.MysqlRowsCloseErr)
		}
	}(rows)
	var adminIds bytes.Buffer
	admins := make([]*vo.QueryAdminsData, 0)
	for rows.Next() {
		a := new(vo.QueryAdminsData)
		err = rows.Scan(&a.Id, &a.Name, &a.Phone, &a.CreateTime)
		if err != nil {
			return enp.Put(errorcode.MysqlScanErr, enp.AddError(err))
		}
		admins = append(admins, a)
		adminIds = tools.ConcatWith(a.Id, adminIds, tools.ConcatWithComma)
	}
	arm, resp := screw.QueryRolesByAdminIds(adminIds.String())
	if resp.Code != errorcode.Success {
		return resp
	}
	for _, admin := range admins {
		if arm[admin.Id] == nil {
			continue
		}
		admin.RoleName = arm[admin.Id].Name
	}
	response.Admins = admins
	return enp.Put(errorcode.Success, enp.AddData(response))
}

func AddAdmin(data []byte) (response *enp.Response) {
	param := new(vo.AddAdminParam)
	err := json.Unmarshal(data, param)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddError(err))
	}
	if len(param.Name) == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("Name", param.Name))
	}
	if len(param.Phone) == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("Phone", param.Name))
	}
	if len(param.Password) == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("Password", param.Name))
	}
	if param.RoleId == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("RoleId", param.Name))
	}
	if len(param.HolderInformation) == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("HolderInformation", param.Name))
	}
	logWriterParam := &log.AddAdminLogWriterParam{Param: param}

	role, resp := L2.GetRoleById(param.RoleId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if role.IsDeleted {
		return enp.Put(errorcode.RoleIsDeleted)
	}
	// 电话号码重复检查
	adminId, resp := L2.GetAdminByPhone(param.Phone)
	if resp.Code != errorcode.Success {
		return resp
	}
	if adminId != 0 {
		return enp.Put(errorcode.AdminRepeated, enp.AddIn(param.Phone))
	}
	tx, err := config.Info().MysqlClient.Begin()
	if err != nil {
		return enp.Put(errorcode.MysqlTxErr, enp.AddError(err))
	}
	defer func() {
		if response != nil && response.Code == errorcode.Success {
			err = tx.Commit()
			if err != nil {
				enp.Put(errorcode.MysqlCommit, enp.AddError(err))
			}
		} else {
			err = tx.Rollback()
			enp.Put(errorcode.MysqlRollback, enp.AddError(err))
		}
	}()
	result, err := tx.Exec("INSERT INTO `trip_portal`.`admin` ( `username`, `phone`, `password`, `create_time`) VALUES (?,?,?,?)",
		param.Name, param.Phone, param.Password, time.Now().Unix())
	if err != nil {
		return enp.Put(errorcode.MysqlExecErr, enp.AddError(err))
	}
	id, err := result.LastInsertId()
	if err != nil {
		return enp.Put(errorcode.MysqlLastInsertIdErr, enp.AddError(err))
	}
	// 管理员与角色的映射关系
	_, err = tx.Exec("INSERT INTO `trip_portal`.`admin_roles` (`admin_id`, `role_id`) VALUES (?, ?)", id, param.RoleId)
	if err != nil {
		return enp.Put(errorcode.MysqlExecErr, enp.AddError(err))
	}
	L2.DelAdminById(id)
	L2.DelAdminByPhone(param.Phone)
	L2.DelAdminRole(id)
	holder := new(do.Admin)
	err = json.Unmarshal([]byte(param.HolderInformation), holder)
	if err != nil {
		return enp.Put(errorcode.MysqlExecErr, enp.AddIn(param), enp.AddError(err))
	}
	logWriterParam.Holder = holder
	log.Add(logWriterParam)
	return enp.Put(errorcode.Success)
}

func GetAdminById(data []byte) *enp.Response {
	param := new(vo.GetAdminByIdParam)
	err := json.Unmarshal(data, param)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddError(err))
	}
	if param.AdminId == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("AdminId", param.AdminId))
	}
	admin, resp := L2.GetAdminById(param.AdminId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if admin == nil || admin.Id == 0 {
		return enp.Put(errorcode.GetAdminByIdNil, enp.AddOut(admin))
	}
	roleId, resp := L2.GetAdminRole(admin.Id)
	if resp.Code != errorcode.Success {
		return resp
	}
	if roleId == 0 {
		return enp.Put(errorcode.GetAdminRoleNil, enp.AddIn(admin.Id))
	}
	adminResp := vo.GetAdminByIdResponse{Id: admin.Id, Name: admin.UserName, Phone: admin.Phone, RoleId: roleId}
	return enp.Put(errorcode.Success, enp.AddData(adminResp))
}

func UpdateAdmin(data []byte) (response *enp.Response) {
	param := new(vo.UpdateAdminParam)
	err := json.Unmarshal(data, param)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddError(err))
	}
	if param.AdminId == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("AdminId", param.AdminId))
	}
	if len(param.HolderInformation) == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("HolderInformation", param.HolderInformation))
	}
	var logWriterParam = &log.UpdateAdminLogWriterParam{Param: param}
	// 获取被修改管理员当前信息
	admin, resp := L2.GetAdminById(param.AdminId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if admin == nil || admin.Id == 0 {
		return enp.Put(errorcode.GetAdminByIdNil, enp.AddIn(param.AdminId))
	}
	logWriterParam.Admin = admin
	// 电话号码重复检查
	adminId, resp := L2.GetAdminByPhone(param.Phone)
	if resp.Code != errorcode.Success {
		return resp
	}
	if adminId != 0 && adminId != param.AdminId {
		return enp.Put(errorcode.AdminRepeated, enp.AddIn(param.Phone))
	}
	// 操作者信息
	holder := new(do.Admin)
	err = json.Unmarshal([]byte(param.HolderInformation), holder)
	if err != nil {
		return enp.Put(errorcode.MysqlExecErr, enp.AddIn(param), enp.AddError(err))
	}
	if holder.Id == 0 {
		return enp.Put(errorcode.AdminHolderInformationNil, enp.AddIn(param))
	}
	logWriterParam.Holder = holder
	tx, err := config.Info().MysqlClient.Begin()
	if err != nil {
		return enp.Put(errorcode.MysqlTxErr, enp.AddError(err))
	}
	defer func() {
		if response != nil && response.Code == errorcode.Success {
			err = tx.Commit()
			if err != nil {
				enp.Put(errorcode.MysqlCommit, enp.AddError(err))
			}
		} else {
			err = tx.Rollback()
			enp.Put(errorcode.MysqlRollback, enp.AddError(err))
		}
	}()
	roleId, resp := L2.GetAdminRole(param.AdminId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if roleId == 0 {
		return enp.Put(errorcode.GetAdminRoleNil, enp.AddIn(param.AdminId))
	}
	if roleId != param.RoleId {
		// 旧权限数据
		oldRole, resp := L2.GetRoleById(roleId)
		if resp.Code != errorcode.Success {
			return resp
		}
		if oldRole == nil || oldRole.Id == 0 {
			return enp.Put(errorcode.GetRoleByIdNil, enp.AddIn(roleId))
		}
		logWriterParam.OldRole = oldRole
		// 新权限数据
		role, resp := L2.GetRoleById(param.RoleId)
		if resp.Code != errorcode.Success {
			return resp
		}
		if role == nil || role.Id == 0 {
			return enp.Put(errorcode.GetRoleByIdNil, enp.AddIn(param.RoleId))
		}
		logWriterParam.Role = role
		// 管理员与角色的映射关系
		// 删除旧数据
		_, err = tx.Exec("DELETE FROM `trip_portal`.`admin_roles` WHERE admin_id=?", param.AdminId)
		if err != nil {
			return enp.Put(errorcode.MysqlExecErr, enp.AddError(err))
		}
		// 增加新数据
		_, err = tx.Exec("INSERT INTO `trip_portal`.`admin_roles` (`admin_id`, `role_id`) VALUES (?, ?)", param.AdminId, param.RoleId)
		if err != nil {
			return enp.Put(errorcode.MysqlExecErr, enp.AddError(err))
		}
		L2.DelAdminRole(param.AdminId)
	}
	var sqlBuf, subSqlBuf bytes.Buffer
	sqlBuf.WriteString("UPDATE `trip_portal`.`admin` SET ")
	var sqlParams = make([]any, 0)
	if len(param.Name) != 0 && param.Name != admin.UserName {
		subSqlBuf = tools.ConcatWith("`username`=?", subSqlBuf, tools.ConcatWithComma)
		sqlParams = append(sqlParams, param.Name)
	}
	if len(param.Phone) != 0 && param.Phone != admin.Phone {
		subSqlBuf = tools.ConcatWith("`phone`=?", subSqlBuf, tools.ConcatWithComma)
		sqlParams = append(sqlParams, param.Phone)
	}
	if len(param.Password) != 0 && param.Password != admin.Password {
		subSqlBuf = tools.ConcatWith("`password`=?", subSqlBuf, tools.ConcatWithComma)
		sqlParams = append(sqlParams, param.Password)
	}
	if subSqlBuf.Len() > 0 {
		subSqlBuf = tools.ConcatWith("`update_time`=?", subSqlBuf, tools.ConcatWithComma)
		sqlParams = append(sqlParams, time.Now().Unix())
		sqlBuf.WriteString(subSqlBuf.String())
		sqlBuf.WriteString(" WHERE `id`=?")
		sqlParams = append(sqlParams, param.AdminId)
		_, err = tx.Exec(sqlBuf.String(), sqlParams...)
		if err != nil {
			return enp.Put(errorcode.MysqlExecErr, enp.AddError(err))
		}
	}
	L2.DelAdminById(param.AdminId)
	L2.DelAdminByPhone(param.Phone)
	L2.DelAdminRole(param.AdminId)
	// 写日志
	log.Add(logWriterParam)
	return enp.Put(errorcode.Success)
}

func DeleteAdmin(data []byte) (response *enp.Response) {
	param := new(vo.UpdateAdminParam)
	err := json.Unmarshal(data, param)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddError(err))
	}
	if param.AdminId == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("AdminId", param.AdminId))
	}
	if len(param.HolderInformation) == 0 {
		return enp.Put(errorcode.InvalidParam, enp.FormatMsg("HolderInformation", param.Name))
	}
	holder := new(do.Admin)
	err = json.Unmarshal([]byte(param.HolderInformation), holder)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddError(err))
	}
	if holder.Id == 0 {
		return enp.Put(errorcode.AdminHolderInformationNil)
	}
	if holder.Id == param.AdminId {
		return enp.Put(errorcode.AdminDeleteItself, enp.AddIn(param.HolderInformation))
	}
	var logWriterParam = &log.DeleteAdminLogWriterParam{Holder: holder}
	// 获取被删除的管理员信息
	admin, resp := L2.GetAdminById(param.AdminId)
	if resp.Code != errorcode.Success {
		return resp
	}
	if admin == nil || admin.Id == 0 {
		return enp.Put(errorcode.GetAdminByIdNil, enp.AddIn(param.AdminId))
	}
	logWriterParam.Admin = admin
	tx, err := config.Info().MysqlClient.Begin()
	if err != nil {
		return enp.Put(errorcode.MysqlTxErr, enp.AddError(err))
	}
	defer func() {
		if response != nil && response.Code == errorcode.Success {
			err = tx.Commit()
			if err != nil {
				enp.Put(errorcode.MysqlCommit, enp.AddError(err))
			}
		} else {
			err = tx.Rollback()
			enp.Put(errorcode.MysqlRollback, enp.AddError(err))
		}
	}()
	// 删除管理员信息
	_, err = tx.Exec("UPDATE `trip_portal`.`admin` SET  `is_delete`=true,`update_time`=? WHERE `id`=?",
		time.Now().Unix(), param.AdminId)
	if err != nil {
		return enp.Put(errorcode.MysqlExecErr, enp.AddError(err))
	}
	// 删除 admin role 关联信息
	_, err = tx.Exec("DELETE FROM `trip_portal`.`admin_roles` WHERE admin_id=?", param.AdminId)
	if err != nil {
		return enp.Put(errorcode.MysqlExecErr, enp.AddError(err))
	}
	err = json.Unmarshal([]byte(param.HolderInformation), holder)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddIn(param), enp.AddError(err))
	}
	log.Add(logWriterParam)
	return enp.Put(errorcode.Success)
}
