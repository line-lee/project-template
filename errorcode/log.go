package errorcode

const (
	LogTimeYear = log*10000 + 1
)

func LogCode() {
	ErrorCode[LogTimeYear] = "开始时间和结束时间只能在同一年份"
}
