package log

import (
	"encoding/json"
	"fmt"
	do2 "github.com/project-template/common/models/variety/do"
	vo2 "github.com/project-template/common/models/variety/vo"
)

const (
	roleAllLog    = roleLog + 0 // 特别注意：二级【全部】选项，必须等于一级选项，（意思就是这里必须【+0】）在处理日志类型搜索时做了相关计算
	addRoleLog    = roleLog + 1
	updateRoleLog = roleLog + 2
	deleteRoleLog = roleLog + 3
)

func roleLogType() {
	set(vo2.LogType{Id: roleLog, Name: "角色", Children: []vo2.LogType{
		{Id: roleAllLog, Name: "全部"},
		{Id: addRoleLog, Name: "新增"},
		{Id: updateRoleLog, Name: "修改"},
		{Id: updateRoleLog, Name: "删除"},
	}})
}

type AddRoleLogWriterParam struct {
	Param *vo2.AddRoleParam
	Role  *do2.Role
}

func (lwp *AddRoleLogWriterParam) writer() *do2.Log {
	holder := new(do2.Admin)
	_ = json.Unmarshal([]byte(lwp.Param.HolderInformation), holder)

	memo := fmt.Sprintf("管理员【%s】【%s】，新增角色【%s】", holder.UserName, holder.Phone, lwp.Role.Name)
	b, _ := json.Marshal(lwp)
	return &do2.Log{Type: roleLog, TypeSub: addRoleLog, AdminId: holder.Id, Memo: memo, IP: holder.RealIP, LogDetail: b}
}
