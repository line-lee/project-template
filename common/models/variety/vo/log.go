package vo

type QueryLogsParam struct {
	TimeStart int64 `json:"time_start" form:"time_start"`
	TimeEnd   int64 `json:"time_end" form:"time_end"`
	AdminId   int64 `json:"admin_id" form:"admin_id"`
	LogType   int64 `json:"log_type" form:"log_type"`
	Page      int64 `json:"page" form:"page"`
	Limit     int64 `json:"limit" form:"limit"`
}

type QueryLogsResponse struct {
	Total int64            `json:"total"`
	Logs  []*QueryLogsData `json:"logs"`
}

type QueryLogsData struct {
	Memo string `json:"memo"`
	IP   string `json:"ip"`
	Time string `json:"time"`
}

type QueryLogTypesResponse struct {
	Types []LogType `json:"types"`
}

type LogType struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Children []LogType `json:"children,omitempty"`
}
