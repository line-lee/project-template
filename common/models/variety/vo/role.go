package vo

type QueryAllAuthsParam struct {
	HolderInformation string
}

type QueryAllAuthsResponse struct {
	MenuStr string `json:"menu_str"`
}

type AddRoleParam struct {
	Name        string `json:"name" form:"name"`
	MenuStr     string `json:"menu_str" form:"menu_str"`
	ButtonStr   string `json:"button_str" form:"button_str"`
	PageStr     string `json:"page_str" form:"page_str"`
	Description string `json:"description" form:"description"`

	HolderInformation string
}
