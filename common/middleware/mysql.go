package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/errorcode"
	"os"
	"runtime"
	"strings"
	"time"
)

func OpenMysqlConnect(cfg *config.Config) {
	cfg.MysqlClient = connect(cfg.MysqlConfig.Username, cfg.MysqlConfig.Password, cfg.MysqlConfig.Host, cfg.MysqlConfig.Port)
	cfg.TripMysqlClientA = connect(cfg.TripMysqlConfigA.Username, cfg.TripMysqlConfigA.Password, cfg.TripMysqlConfigA.Host, cfg.TripMysqlConfigA.Port)
	cfg.TripMysqlClientB = connect(cfg.TripMysqlConfigB.Username, cfg.TripMysqlConfigB.Password, cfg.TripMysqlConfigB.Host, cfg.TripMysqlConfigB.Port)
	cfg.TripMysqlClientC = connect(cfg.TripMysqlConfigC.Username, cfg.TripMysqlConfigC.Password, cfg.TripMysqlConfigC.Host, cfg.TripMysqlConfigC.Port)
}

func connect(username, password, host string, port int) *sql.DB {
	url := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8", username, password, host, port)
	db, err := sql.Open("mysql", url)
	if err != nil {
		fmt.Println("mysql open err ", err)
		os.Exit(-1)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("mysql ping err ", err)
		os.Exit(-1)
	}
	db.SetMaxOpenConns(runtime.NumCPU()<<1 + 2)
	db.SetMaxIdleConns(runtime.NumCPU()<<1 + 2)
	db.SetConnMaxLifetime(3600 * time.Second)
	return db
}

func CloseMysqlConnect() {
	err := config.Info().MysqlClient.Close()
	if err != nil {
		fmt.Println("mysql connect close err ", err)
		os.Exit(-1)
	}
	err = config.Info().TripMysqlClientA.Close()
	if err != nil {
		fmt.Println("trip mysql a connect close err ", err)
		os.Exit(-1)
	}
	err = config.Info().TripMysqlClientB.Close()
	if err != nil {
		fmt.Println("trip mysql b connect close err ", err)
		os.Exit(-1)
	}
	err = config.Info().TripMysqlClientC.Close()
	if err != nil {
		fmt.Println("trip mysql c connect close err ", err)
		os.Exit(-1)
	}
}

// 内存中查询存在的数据库表
var tm = make(map[string]bool)

const (
	DBHour  = 1
	DBDay   = 2
	DBMonth = 3
	DBYear  = 4
)

const (
	TableLog           = "log"             // 操作日志表
	TableLogDetail     = "log_detail"      // 操作日志变更参数记录表
	TableTripAccessLog = "trip_access_log" // 与网约车系统连接错误日志记录表
)

func IsDBExist(tableName string, thisTime time.Time, timeRange int) (string, *enp.Response) {
	// 循环检查，cm val 赋值表名
	timeTage, lastTimeTags, resp := timeTag(thisTime, timeRange)
	if resp.Code != errorcode.Success {
		return "", resp
	}
	// 使用内存减少数据库查询
	var expectTableName, lastTableName string
	expectTableName = fmt.Sprintf("%s%s", tableName, timeTage)

	if isExist := tm[expectTableName]; !isExist {
		// 内存不存在，继续查库
		tableCheckSql := "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'trip_portal' AND TABLE_NAME = ?"
		var tn string
		err := config.Info().MysqlClient.QueryRow(tableCheckSql, expectTableName).Scan(&tn)
		if len(tn) != 0 {
			// 表存在，写入内存，返回正确
			tm[expectTableName] = true
			return expectTableName, enp.Put(errorcode.Success)
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return "", enp.Put(errorcode.MysqlScanErr, enp.AddError(err))
		}
		for _, lastTimeTag := range lastTimeTags {
			lastTableName = fmt.Sprintf("%s%s", tableName, lastTimeTag)
			if resp = createTable(lastTableName, expectTableName); resp.Code == errorcode.Success {
				tm[expectTableName] = true
				return expectTableName, enp.Put(errorcode.Success)
			}
			// 建表错误，继续循环执行
			// .....
		}
		// 建表循环结束，都没成功
		return "", enp.Put(errorcode.MysqlSharding, enp.AddIn(lastTimeTags))
	}
	return expectTableName, enp.Put(errorcode.Success)
}

func timeTag(thisTime time.Time, timeRange int) (string, []string, *enp.Response) {
	var thisTimeTage = thisTime.Format("2006010215")
	var lastTimeTags = make([]string, 0)
	for i := 1; i <= 10; i++ {
		switch timeRange {
		case DBHour:
			thisTimeTage = thisTime.Format("2006010215")
			lastTime := thisTime.Add(-time.Duration(i) * time.Hour)
			lastTimeTags = append(lastTimeTags, lastTime.Format("2006010215"))
		case DBDay:
			thisTimeTage = thisTime.Format("20060102")
			lastTime := thisTime.AddDate(0, 0, -i)
			lastTimeTags = append(lastTimeTags, lastTime.Format("20060102"))
		case DBMonth:
			thisTimeTage = thisTime.Format("200601")
			lastTime := thisTime.AddDate(0, -i, 0)
			lastTimeTags = append(lastTimeTags, lastTime.Format("200601"))
		case DBYear:
			thisTimeTage = thisTime.Format("2006")
			lastTime := thisTime.AddDate(-i, 0, 0)
			lastTimeTags = append(lastTimeTags, lastTime.Format("2006"))
		default:
			return "", nil, enp.Put(errorcode.MysqlShardingTimeRangeUnknown)
		}
	}
	return thisTimeTage, lastTimeTags, enp.Put(errorcode.Success)
}

func createTable(lastTableName, expectTableName string) *enp.Response {
	// 表不存在，copy上一个建表结构，新建表
	showCreateSql := fmt.Sprintf("SHOW CREATE TABLE trip_portal.%s", lastTableName)
	var showTableName, createSql string
	err := config.Info().MysqlClient.QueryRow(showCreateSql).Scan(&showTableName, &createSql)
	if err != nil {
		return enp.Put(errorcode.MysqlScanErr, enp.AddError(err))
	}
	createSql = strings.ReplaceAll(createSql, fmt.Sprintf("`%s`", lastTableName), fmt.Sprintf("`trip_portal`.`%s`", expectTableName))
	_, err = config.Info().MysqlClient.Exec(createSql)
	if err != nil {
		return enp.Put(errorcode.MysqlExecErr, enp.AddError(err))
	}
	return enp.Put(errorcode.Success)
}
