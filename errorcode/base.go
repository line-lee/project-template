package errorcode

var ErrorCode = make(map[int]string)

const (
	base       = 1
	sdk        = 2
	middleware = 3
	trip       = 4
	admin      = 5
	role       = 6
	merchant   = 7
	log        = 8
)

const (
	Success        = base*10000 + 1
	Recover        = base*10000 + 2
	SystemError    = base*10000 + 3
	Unauthorized   = base*10000 + 4
	InvalidParam   = base*10000 + 5
	UnknownModule  = base*10000 + 6
	UnknownService = base*10000 + 7
)

func BaseCode() {
	ErrorCode[Success] = "Success"
	ErrorCode[Recover] = "Recover"
	ErrorCode[SystemError] = "系统错误"
	ErrorCode[Unauthorized] = "鉴权未通过，请重新登录#?#鉴权失败具体原因：[%v]"
	ErrorCode[InvalidParam] = "无效参数#?#参数检查未通过[%v][%v]"
	ErrorCode[UnknownModule] = "grpc 下游解析 module 不识别"
	ErrorCode[UnknownService] = "grpc 下游解析 service 不识别"
}
