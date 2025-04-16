package errorcode

const (
	GetRoleByIdNil       = role*10000 + 1
	QueryRolesError      = role*1000 + 2
	AddRoleAlreadyExist  = role*1000 + 3
	RoleIsDeleted        = role*1000 + 4
	RolePass             = role*1000 + 5
	AdminRoleNil         = role*1000 + 6
	AdminAlreadyBindRole = role*1000 + 7
)

func RoleCode() {
	ErrorCode[GetRoleByIdNil] = "根据id，查询角色信息为空"
	ErrorCode[QueryRolesError] = "查询角色列表信息错误"
	ErrorCode[AddRoleAlreadyExist] = "新增角色已存在"
	ErrorCode[RoleIsDeleted] = "所选角色权限已被删除，请重新选择"
	ErrorCode[RolePass] = "无权操作"
	ErrorCode[AdminRoleNil] = "管理员绑定角色为空"
	ErrorCode[AdminAlreadyBindRole] = "该角色已被成员账号绑定，无法删除"
}
