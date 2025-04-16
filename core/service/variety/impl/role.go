package impl

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/project-template/cache/L2"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/do"
	"github.com/project-template/common/models/variety/vo"
	"github.com/project-template/errorcode"
	"github.com/project-template/screw/log"
	"time"
)

// QueryAllAuths TODO FIX
func QueryAllAuths(data []byte) *enp.Response {
	menus, resp := queryMenus()
	if resp.Code != errorcode.Success {
		return resp
	}
	if menus != nil {
		for i := 0; i < len(menus); i++ {
			newPages := make([]*do.Page, 0)
			pages, resp := queryPages(menus[i].Id)
			if resp.Code != errorcode.Success {
				return resp
			}
			if pages != nil {
				for i := 0; i < len(pages); i++ {
					newPages = append(newPages, pages[i])
					buttons, resp := queryButtons(pages[i].Id)
					if resp.Code != errorcode.Success {
						return resp
					}
					pages[i].Buttons = buttons
				}
			}
			menus[i].Pages = newPages
		}
	}
	var menuStr string
	if len(menus) > 0 {
		bytes, _ := json.Marshal(menus)
		menuStr = string(bytes)
	}
	return enp.Put(errorcode.Success, enp.AddData(vo.QueryAllAuthsResponse{MenuStr: menuStr}))
}

func queryMenus() ([]*do.Menu, *enp.Response) {
	rows, err := config.Info().MysqlClient.Query("select `id`,`name`,`path`,`sort` from `trip_portal`.`menu`")
	if err != nil {
		return nil, enp.Put(errorcode.MysqlQueryErr, enp.AddError(err))
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			enp.Put(errorcode.MysqlRowsCloseErr, enp.AddError(err))
		}
	}(rows)

	var list []*do.Menu
	for rows.Next() {
		menu := new(do.Menu)
		err = rows.Scan(&menu.Id, &menu.Name, &menu.Path, &menu.Sort)
		if err != nil {
			return nil, enp.Put(errorcode.MysqlScanErr, enp.AddError(err))
		}
		list = append(list, menu)
	}
	return list, enp.Put(errorcode.Success)
}

func queryPages(menuId int64) ([]*do.Page, *enp.Response) {
	rows, err := config.Info().MysqlClient.Query("select `id`,`name`,`menu_id`,`sort` from `trip_portal`.`page` where `menu_id` = ?", menuId)
	if err != nil {
		return nil, enp.Put(errorcode.MysqlQueryErr, enp.AddError(err))
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			enp.Put(errorcode.MysqlRowsCloseErr, enp.AddError(err))
		}
	}(rows)

	var list []*do.Page
	for rows.Next() {
		page := new(do.Page)
		err = rows.Scan(&page.Id, &page.Name, &page.MenuId, &page.Sort)
		if err != nil {
			return nil, enp.Put(errorcode.MysqlScanErr, enp.AddError(err))
		}
		list = append(list, page)
	}
	return list, enp.Put(errorcode.Success)
}

func queryButtons(pageId int64) ([]*do.Button, *enp.Response) {
	rows, err := config.Info().MysqlClient.Query("select `id`,`name`,`page_id`,`sort` from `trip_portal`.`button` where `page_id` = ?", pageId)
	if err != nil {
		return nil, enp.Put(errorcode.MysqlQueryErr, enp.AddError(err))
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			enp.Put(errorcode.MysqlRowsCloseErr, enp.AddError(err))
		}
	}(rows)
	var list []*do.Button
	for rows.Next() {
		button := new(do.Button)
		err = rows.Scan(&button.Id, &button.Name, &button.PageId, &button.Sort)
		if err != nil {
			return nil, enp.Put(errorcode.MysqlScanErr, enp.AddError(err))
		}
		list = append(list, button)
	}
	return list, enp.Put(errorcode.Success)
}

func AddRole(data []byte) *enp.Response {
	req := new(vo.AddRoleParam)
	err := json.Unmarshal(data, req)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal)
	}
	var total int32
	err = config.Info().MysqlClient.QueryRow("select count(id) from `trip_portal`.`roles` where `name` = ? and `is_deleted` = false ", req.Name).Scan(&total)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return enp.Put(errorcode.SystemError)
	}
	if total > 0 {
		return enp.Put(errorcode.AddRoleAlreadyExist)
	}
	result, err := config.Info().MysqlClient.Exec("insert into `trip_portal`.`roles` (`name`,`description`,`menu_str`,`button_str`,`page_str`,`is_main`,`is_deleted`,`updated`,`created`,`version`) value(?,?,?,?,?, ?,?,?,?,?)",
		req.Name, req.Description, req.MenuStr, req.ButtonStr, req.PageStr, false, false, time.Now().Unix(), time.Now().Unix(), 1)
	if err != nil {
		return enp.Put(errorcode.MysqlExecErr, enp.AddError(err))
	}
	id, err := result.LastInsertId()
	if err != nil {
		return enp.Put(errorcode.MysqlLastInsertIdErr, enp.AddError(err))
	}
	L2.DelRole(id)
	log.Add(&log.AddRoleLogWriterParam{Param: req, Role: &do.Role{Name: req.Name}})
	return enp.Put(errorcode.Success)
}
