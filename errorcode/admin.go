package errorcode

const (
	GetAdminByIdNil           = admin*10000 + 1
	AdminTokenExpire          = admin*10000 + 2
	GetAdminRoleNil           = admin*10000 + 3
	GetAdminByPhoneNil        = admin*10000 + 4
	AdminDelete               = admin*10000 + 5
	AdminLock                 = admin*10000 + 6
	AdminPassword             = admin*10000 + 7
	AdminRepeated             = admin*10000 + 8
	AdminHolderInformationNil = admin*10000 + 9
	AdminDeleteItself         = admin*10000 + 10
	AdminTokenSSO             = admin*10000 + 11
	AdminTokenIdNil           = admin*10000 + 12
)

func AdminCode() {
	ErrorCode[GetAdminByIdNil] = "根据id，查询admin信息为空"
	ErrorCode[AdminTokenExpire] = "token expire"
	ErrorCode[GetAdminRoleNil] = "根据admin id，查询admin role 信息为空"
	ErrorCode[GetAdminByPhoneNil] = "账号不存在#?#根据账号phone查询管理员信息为空"
	ErrorCode[AdminDelete] = "账号不存在#?#管理员账号已被删除"
	ErrorCode[AdminLock] = "账号不存在#?#管理员账号已被锁定"
	ErrorCode[AdminPassword] = "账号或密码不正确#?#管理员账号密码不匹配"
	ErrorCode[AdminRepeated] = "电话号码已注册管理员"
	ErrorCode[AdminDeleteItself] = "不允许删除自己的账号"
	ErrorCode[AdminTokenSSO] = "token sso"
	ErrorCode[AdminTokenIdNil] = "token id nil"
}
