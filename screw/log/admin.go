package log

import (
	"encoding/json"
	"fmt"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/common/models/variety/vo"
)

const (
	adminAllLog    = adminLog + 0 // 特别注意：二级【全部】选项，必须等于一级选项（意思就是这里必须【+0】），在处理日志类型搜索时做了相关计算
	adminLoginLog  = adminLog + 1
	addAdminLog    = adminLog + 2
	updateAdminLog = adminLog + 3
	deleteAdminLog = adminLog + 4
)

func adminLogType() {
	set(vo.LogType{Id: adminLog, Name: "管理员", Children: []vo.LogType{
		{Id: adminAllLog, Name: "全部"},
		{Id: adminLoginLog, Name: "登录"},
		{Id: addAdminLog, Name: "新增"},
		{Id: updateAdminLog, Name: "修改"},
		{Id: deleteAdminLog, Name: "删除"},
	}})
}

type AdminLoginLogWriterParam struct {
	Param *vo.AdminLoginParam // 管理员登录参数
	Admin *do.Admin           // 管理员信息
}

func (lwp *AdminLoginLogWriterParam) writer() *do.Log {
	var param = lwp.Param
	var admin = lwp.Admin
	memo := fmt.Sprintf("管理员【%s】【%s】，登录系统", admin.UserName, admin.Phone)
	b, _ := json.Marshal(lwp)
	return &do.Log{Type: adminLog, TypeSub: adminLoginLog, AdminId: admin.Id, Memo: memo, IP: param.RealIP, LogDetail: b}
}
