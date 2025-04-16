package vo

type AdminLoginParam struct {
	Phone    string `json:"phone" form:"phone"`       //电话
	Password string `json:"password" form:"password"` // 密码

	// 不需要传入，自己解析
	RealIP string
}

type AdminLoginResponse struct {
	Name         string `json:"name"`          // 账号名称
	Phone        string `json:"phone"`         // 账号电话
	Menu         string `json:"menu"`          // 菜单权限集合
	Page         string `json:"page"`          // 页面权限集合
	Button       string `json:"button"`        //按钮权限集合
	Token        string `json:"token"`         // token
	RefreshToken string `json:"refresh_token"` // 延时token
}

type AdminTokenRefreshParam struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type AdminTokenRefreshResponse struct {
	Name         string `json:"name"`          // 账号名称
	Phone        string `json:"phone"`         // 账号电话
	Menu         string `json:"menu"`          // 菜单权限集合
	Page         string `json:"page"`          // 页面权限集合
	Button       string `json:"button"`        //按钮权限集合
	Token        string `json:"token"`         // token
	RefreshToken string `json:"refresh_token"` // 延时token
}

type QueryAdminsParam struct {
	Phone string `json:"phone" form:"phone"` // 管理员电话号码
	Name  string `json:"name" form:"name"`
	Page  int    `json:"page" form:"page"`
	Limit int    `json:"limit" form:"limit"`
}

type QueryAdminsResponse struct {
	Total  int64              `json:"total"`
	Admins []*QueryAdminsData `json:"admins,omitempty"`
}

type QueryAdminsData struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`        // 管理员账号名称
	Phone      string `json:"phone"`       // 管理员电话
	RoleName   string `json:"role_name"`   // 角色名称
	CreateTime int64  `json:"create_time"` // 创建时间
}

type AddAdminParam struct {
	Name     string `json:"name" form:"name"`         // 管理员账号名称
	Phone    string `json:"phone" form:"phone"`       // 管理员电话
	Password string `json:"password" form:"password"` // 密码
	RoleId   int64  `json:"role_id" form:"role_id"`   // 角色主键

	HolderInformation string `json:"holder_information"`
}

type GetAdminByIdParam struct {
	AdminId int64 `json:"admin_id" form:"admin_id"`
}

type GetAdminByIdResponse struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`  // 管理员账号名称
	Phone    string `json:"phone"` // 管理员电话
	Password string `json:"password"`
	RoleId   int64  `json:"role_id"` // 角色主键
}

type UpdateAdminParam struct {
	AdminId  int64  `json:"admin_id" form:"admin_id"` // 修改的管理员id
	Name     string `json:"name" form:"name"`         // 管理员账号名称
	Phone    string `json:"phone" form:"phone"`       // 管理员电话
	Password string `json:"password" form:"password"` // 密码
	RoleId   int64  `json:"role_id" form:"role_id"`   // 角色主键

	HolderInformation string `json:"holder_information"`
}

type DeleteAdminParam struct {
	AdminId int64 `json:"admin_id" form:"admin_id"` // 需要删除的管理员id

	HolderInformation string `json:"holder_information"`
}

type UpdateAdminPasswordParam struct {
	AdminId  int64  `json:"admin_id" form:"admin_id"` // 修改其他管理员密码需要传入id，修改自己的不用
	Password string `json:"password" form:"password"`

	HolderInformation string `json:"holder_information"`
}
