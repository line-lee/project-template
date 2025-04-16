package do

type Admin struct {
	Id         int64
	UserName   string // 账户名称
	Phone      string // 电话
	Password   string // 密码
	IsDelete   bool   // 删除标识
	IsLock     bool   // 是否锁定
	CreateTime int64
	UpdateTime int64

	// 冗余字段，不落库，只做逻辑计算
	RealIP string // 管理员登录真实ip
}
