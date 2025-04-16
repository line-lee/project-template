package do

type Menu struct {
	Id    int64   `json:"id"`
	Name  string  `json:"name"`
	Path  string  `json:"path"`
	Sort  int32   `json:"sort"`
	Pages []*Page `json:"pages"`
}

type Page struct {
	Id      int64     `json:"id"`
	Name    string    `json:"name"`
	MenuId  int64     `json:"menu_id"`
	Sort    int32     `json:"sort"`
	Buttons []*Button `json:"buttons"`
}

type Button struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Sort   int32  `json:"sort"`
	PageId int64  `json:"page_id"`
}

type Role struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`        // 角色名称
	Description string `json:"description"` // 描述
	MenuStr     string `json:"menu_str"`    // 菜单
	PageStr     string `json:"page_str"`    // 页面
	ButtonStr   string `json:"button_str"`  // 按钮
	IsDeleted   bool   `json:"is_deleted"`
	IsMain      bool   `json:"is_main"` // 是否为超级管理员
	Updated     int64  `json:"updated"`
	Created     int64  `json:"created"`
	Version     int32  `json:"version"`
}

type AdminRole struct {
	AdminId int64 // 管理员id
	RoleId  int64 // 角色id
}
