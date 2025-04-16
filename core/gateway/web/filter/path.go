package filter

// 全部filter免检查
var completeCoverage = map[string]bool{
	"/web/api/v1/login": true,
}

// 免除token过期检查
var tokenExpire = map[string]bool{
	"/web/api/v1/token/refresh": true,
}
