package bo

import "time"

const (
	HolderInformation  = "holder_information" // 接口信息持有者，上下文传递ctx
	HolderInformationX = "HolderInformation"  //接口信息持有者，上下文传递ctx

	TokenShort = 1 // 短有效期
	TokenLong  = 2 // 长有效期

	AdminTokenKey         = "x-admin-holder" // 与web交互，放在 request header 的 token key
	AdminContextKey       = "admin_claim"    // gateway 各个 filter 中间件 执行传递的对象 上下文 key
	AdminTokenShortExpire = 30 * time.Minute // 短 token 有效期 30 分钟
	AdminTokenLongExpire  = 8 * time.Hour    // 长 token 有效期 8 小时
)

type AdminClaim struct {
	Id         int64
	UserName   string // 名称
	Phone      string // 电话
	Password   string // 密码
	Type       int64  // 1.短有效期；长有效期
	SSO        string // 单点登录校验标识（single sign on）
	ExpireTime int64  // 过期时间
	RoleId     int64  // 角色id
	MenuStr    string // 菜单权限集
	PageStr    string // 页面权限集
	ButtonStr  string // 按钮权限集
}
