package vo

type BaiduLocationResponse struct {
	Provinces []*Province `json:"provinces"`
}

type Province struct {
	Name   string  `json:"name"`
	Cities []*City `json:"cities"`
}

type City struct {
	Code      string      `json:"code"`
	Name      string      `json:"name"`
	Districts []*District `json:"districts"`
}

type District struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type PasswordGenerateResponse struct {
	Password string `json:"password"` // 生成的密码
}

type GetQiNiuTokenResponse struct {
	Token string `json:"token"`
}
