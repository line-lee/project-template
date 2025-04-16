package listener

import (
	"encoding/json"
	"fmt"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	mw "github.com/project-template/common/middleware"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/errorcode"
	"time"
)

func log() {
	// 操作日志
	mw.KafkaConsumer(mw.KafkaTopicLog, func(message []byte) (response *enp.Response) {
		km := new(mw.KafkaMessage)
		err := json.Unmarshal(message, km)
		if err != nil {
			return enp.Put(errorcode.JsonUnmarshal)
		}
		l := new(do.Log)
		err = json.Unmarshal(km.Val, l)
		if err != nil {
			return enp.Put(errorcode.JsonUnmarshal)
		}
		var logTableName, logDetailTableName string
		var resp = new(enp.Response)
		if logTableName, resp = mw.IsDBExist(mw.TableLog, km.Time, mw.DBYear); resp.Code != errorcode.Success {
			return resp
		}
		if logDetailTableName, resp = mw.IsDBExist(mw.TableLogDetail, km.Time, mw.DBYear); resp.Code != errorcode.Success {
			return resp
		}
		var logSql = fmt.Sprintf("INSERT INTO `trip_portal`.`%s` (`admin_id`, `type`,`type_sub`, `memo`, `ip`,`log_time`, `create_time`) VALUES ( ?,?, ?, ?, ?, ?, ?)", logTableName)
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
				if err != nil {
					enp.Put(errorcode.MysqlRollback, enp.AddError(err))
				}
			}
		}()
		result, err := tx.Exec(logSql, l.AdminId, l.Type, l.TypeSub, l.Memo, l.IP, km.Time.Unix(), time.Now().Unix())
		if err != nil {
			return enp.Put(errorcode.MysqlExecErr, enp.AddIn(logSql, logTableName), enp.AddError(err))
		}
		id, err := result.LastInsertId()
		if err != nil {
			return enp.Put(errorcode.MysqlLastInsertIdErr, enp.AddIn(logSql, logTableName), enp.AddError(err))
		}
		// 记录日志参数
		var detailSql = fmt.Sprintf("INSERT INTO `trip_portal`.`%s` ( `log_id`, `param`) VALUES ( ?, ?)", logDetailTableName)
		_, err = tx.Exec(detailSql, id, string(l.LogDetail))
		if err != nil {
			return enp.Put(errorcode.MysqlExecErr, enp.AddIn(detailSql, logDetailTableName), enp.AddError(err))
		}
		return enp.Put(errorcode.Success)
	})

	// 系统错误日志......

	// 操作日志
	mw.KafkaConsumer(mw.KafkaTopicLog, func(message []byte) (response *enp.Response) {
		return enp.Put(errorcode.Success)
	})
}
