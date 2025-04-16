package log

import (
	"encoding/json"
	enp "github.com/project-template/common/encapsulate"
	mw "github.com/project-template/common/middleware"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/common/models/variety/vo"
	"github.com/project-template/errorcode"
	"time"
)

const (
	adminLog = 10000 // 管理员相关操作日志
	roleLog  = 20000 // 角色相关操作日志
)

func Add(wt WriterType) {
	lb, err := json.Marshal(wt.writer())
	if err != nil {
		enp.Put(errorcode.JsonMarshal, enp.AddError(err))
		return
	}
	mw.KafkaProduce(mw.KafkaTopicLog, &mw.KafkaMessage{Val: lb, Time: time.Now()})
}

type WriterType interface {
	writer() *do.Log
}

// 日志类型
var lts = make([]vo.LogType, 0)

func set(t vo.LogType) {
	lts = append(lts, t)
}

func logTypeInit() {
	adminLogType()
	roleLogType()
}

func TypeNames() []vo.LogType {
	if len(lts) == 0 {
		logTypeInit()
	}
	return lts
}
