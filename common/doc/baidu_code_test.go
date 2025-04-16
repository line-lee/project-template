package doc

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"testing"
)

type Province struct {
	Name   string
	Cities []*City
}

type City struct {
	Code      string
	Name      string
	Districts []*District
}

type District struct {
	Name string
	Code string
}

// 下载地址：https://lbs.baidu.com/faq/api?title=webapi/download
// 选择->资源下载->国内城市行政区划代码
// 下载文件的原名称是：weather_district_id .xlsx
func TestBaiduCode(t *testing.T) {
	f, err := excelize.OpenFile("E:\\web_work_space\\trip-portal\\common\\doc\\baidu_code.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := f.GetRows("weather_district_id")
	if err != nil {
		fmt.Println(err)
		return
	}

	var provinces = make([]*Province, 0)
	var province *Province
	var provinceFilterMap = make(map[string]bool)
	var city *City
	var cityFilterMap = make(map[string]bool)

	var district *District

	for rowIdx, row := range rows {
		if rowIdx == 0 {
			// 表头
			continue
		}
		// 1.省份名称；2.城市名称；3.cityCode；4.县区名称；5.adCode
		if !provinceFilterMap[row[1]] {
			if province != nil {
				if province.Cities == nil {
					province.Cities = make([]*City, 0)
				}
				if city != nil && len(city.Districts) > 1 {
					city.Districts = city.Districts[1:]
				}
				province.Cities = append(province.Cities, city)
				city = nil
				provinces = append(provinces, province)
			}
			province = &Province{Name: row[1]}
			provinceFilterMap[row[1]] = true
		}
		if province == nil {
			fmt.Println("新建 province 为空")
			return
		}
		if !cityFilterMap[row[2]] {
			if city != nil {
				if province.Cities == nil {
					province.Cities = make([]*City, 0)
				}
				if len(city.Districts) > 1 {
					city.Districts = city.Districts[1:]
				}
				province.Cities = append(province.Cities, city)
				city = nil
			}
			city = &City{Name: row[2], Code: row[3]}
			cityFilterMap[row[2]] = true
		}
		if city == nil {
			fmt.Println("新建 city 为空")
			return
		}
		if city.Code == row[5] {
			district = &District{Name: "全部", Code: row[5]}
		} else {
			district = &District{Name: row[4], Code: row[5]}
		}
		if len(city.Districts) == 0 {
			city.Districts = make([]*District, 0)
		}
		city.Districts = append(city.Districts, district)
	}

	b, _ := json.Marshal(provinces)
	fmt.Println(string(b))
	if err = os.WriteFile("E:\\web_work_space\\trip-portal\\common\\doc\\baidu_code.json", b, 0644); err != nil {
		fmt.Println(err)
	}
}
