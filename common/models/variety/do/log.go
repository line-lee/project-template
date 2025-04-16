package do

type Log struct {
	Id         int64
	AdminId    int64  // 管理员id
	Type       int64  // 日志类型，比如：管理员，商户，收费配置.....
	TypeSub    int64  // 日志子类型，比如：登录，新增，修改.....
	Memo       string // 日志详情
	IP         string // ip
	LogTime    int64  // 日志在业务中调用产生的时间
	CreateTime int64  // 日志最终落库的时间

	// 冗余参数，只做计算，不落库
	LogDetail []byte
}
