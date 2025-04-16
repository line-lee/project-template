package impl

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/vo"
	"github.com/project-template/errorcode"
	"time"
)

func QueryLogs(data []byte) *enp.Response {
	param := new(vo.QueryLogsParam)
	err := json.Unmarshal(data, param)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddIn(string(data)), enp.AddError(err))
	}
	if param.TimeStart == 0 {
		return enp.Put(errorcode.InvalidParam, enp.AddIn(param), enp.FormatMsg("TimeStart", param.TimeStart))
	}
	if param.TimeEnd == 0 {
		return enp.Put(errorcode.InvalidParam, enp.AddIn(param), enp.FormatMsg("TimeEnd", param.TimeStart))
	}
	if param.Page == 0 {
		param.Page = 1
	}
	if param.Limit < 15 || param.Limit > 60 {
		param.Limit = 15
	}
	start := time.Unix(param.TimeStart, 0)
	end := time.Unix(param.TimeEnd, 0)
	if start.Year() != end.Year() {
		return enp.Put(errorcode.LogTimeYear, enp.AddIn(param))
	}
	var tableName = fmt.Sprintf("`trip_portal`.`log%d`", start.Year())
	var countSql, querySql, commonSql bytes.Buffer
	var sqlParam = make([]any, 0)
	countSql.WriteString(fmt.Sprintf("SELECT COUNT(`id`) FROM %s ", tableName))
	querySql.WriteString(fmt.Sprintf("SELECT `memo`,`ip`,`log_time` FROM %s ", tableName))
	commonSql.WriteString("WHERE `log_time` >= ? AND `log_time` <= ? ")
	sqlParam = append(sqlParam, []any{param.TimeStart, param.TimeEnd}...)
	if param.AdminId != 0 {
		commonSql.WriteString("AND `admin_id`=? ")
		sqlParam = append(sqlParam, param.AdminId)
	}
	if param.LogType != 0 {
		if param.LogType%10000 == 0 {
			// 搜索主类型
			commonSql.WriteString("AND `type`=? ")
			sqlParam = append(sqlParam, param.LogType)
		} else {
			commonSql.WriteString("AND `type_sub`=? ")
			sqlParam = append(sqlParam, param.LogType)
		}
	}
	countSql.WriteString(commonSql.String())
	var total int64
	err = config.Info().MysqlClient.QueryRow(countSql.String(), sqlParam...).Scan(&total)
	if err != nil {
		return enp.Put(errorcode.MysqlScanErr, enp.AddIn(countSql.String(), sqlParam), enp.AddError(err))
	}
	if total == 0 {
		return enp.Put(errorcode.Success, enp.AddData(vo.QueryLogsResponse{Total: total}))
	}
	querySql.WriteString(commonSql.String())
	querySql.WriteString("ORDER BY `log_time` DESC LIMIT ? OFFSET ?")
	sqlParam = append(sqlParam, []any{param.Limit, (param.Page - 1) * param.Limit}...)
	rows, err := config.Info().MysqlClient.Query(querySql.String(), sqlParam...)
	if err != nil {
		return enp.Put(errorcode.MysqlQueryErr, enp.AddIn(querySql.String(), sqlParam), enp.AddError(err))
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			enp.Put(errorcode.MysqlRowsCloseErr, enp.AddIn(querySql.String(), sqlParam), enp.AddError(err))
		}
	}(rows)
	var lds = make([]*vo.QueryLogsData, 0)
	for rows.Next() {
		var ld = new(vo.QueryLogsData)
		err = rows.Scan(&ld.Memo, &ld.IP, &ld.Time)
		if err != nil {
			return enp.Put(errorcode.MysqlScanErr, enp.AddIn(querySql.String(), sqlParam), enp.AddError(err))
		}
		lds = append(lds, ld)
	}
	return enp.Put(errorcode.Success, enp.AddData(vo.QueryLogsResponse{Total: total, Logs: lds}))
}
