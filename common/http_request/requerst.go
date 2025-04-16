package http_request

import (
	"bytes"
	"fmt"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/models/variety/bo"
	"github.com/project-template/errorcode"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RequestParam struct {
	Header map[string]interface{} // 请求头
	Path   []string               // 路径参数
	Query  map[string]interface{} // // 查询参数
	Body   map[string]interface{} // 请求体
	Json   []byte
}

const (
	contentType     = "Content-Type"
	contentTypeJSON = "application/json;charset=UTF-8"
	contentTypeForm = "application/x-www-form-urlencoded"
)

func NewRequest(method, address string, param *RequestParam) ([]byte, *enp.Response) {
	var buf bytes.Buffer
	buf.WriteString("--------------------------------------------------------------------\n")
	buf.WriteString(fmt.Sprintf("请求方式：%s  请求路径：%s\n", method, address))
	defer func() {
		buf.WriteString("--------------------------------------------------------------------\n")
		fmt.Println(buf.String())
	}()
	var reader *strings.Reader
	// 路径参数
	if param.Path != nil {
		for _, v := range param.Path {
			address = fmt.Sprintf("%s/%s", address, v)
		}
		buf.WriteString(fmt.Sprintf("路径参数：%s\n", address))
	}
	// 处理查询参数
	if param.Query != nil {
		uri, err := url.Parse(address)
		if err != nil {
			return nil, enp.Put(errorcode.UrlParse, enp.AddError(err))
		}
		var uvs = make(url.Values)
		for k, v := range param.Query {
			uvs[k] = []string{fmt.Sprint(v)}
		}
		uri.RawQuery = uvs.Encode()
		address = uri.String()
		buf.WriteString(fmt.Sprintf("查询参数：%s\n", address))
	}
	// 处理请求体参数
	//if param.Body != nil {
	switch param.Header[contentType] {
	case contentTypeForm:
		values := url.Values{}
		for k, v := range param.Body {
			values.Set(k, fmt.Sprintf("%v", v))
		}
		buf.WriteString(fmt.Sprintf("处理请求体参数[FROM]：%s\n", values.Encode()))
		reader = strings.NewReader(values.Encode())
	case contentTypeJSON:
		//jsonBytes, err := json.Marshal(param.Body)
		//if err != nil {
		//	return nil, enp.Put(errorcode.JsonMarshal, enp.AddError(err))
		//}
		buf.WriteString(fmt.Sprintf("处理请求体参数[JSON]：%s\n", string(param.Json)))
		reader = strings.NewReader(string(param.Json))
	}
	//}
	var request *http.Request
	var err error
	if reader == nil {
		request, err = http.NewRequest(method, address, nil)
	} else {
		request, err = http.NewRequest(method, address, reader)
	}
	if err != nil {
		return nil, enp.Put(errorcode.HttpNewRequest, enp.AddError(err))
	}
	if param.Header != nil {
		buf.WriteString("请求头[Header]:\n")
		for k, v := range param.Header {
			request.Header.Add(k, fmt.Sprint(v))
			buf.WriteString(fmt.Sprintf("  %s:%s\n", k, fmt.Sprint(v)))
		}
	}
	client := http.Client{Timeout: 3 * time.Second}
	start := time.Now()
	response, err := client.Do(request)
	end := time.Now()
	if err != nil {
		return nil, enp.Put(errorcode.HttpClientDo, enp.AddError(err))
	}
	if response == nil {
		buf.WriteString("返回 response nil\n")
		return nil, enp.Put(errorcode.HttpResponseNil)
	}
	if response.StatusCode != http.StatusOK {
		buf.WriteString(fmt.Sprintf("返回 http code %d\n", response.StatusCode))
		return nil, enp.Put(errorcode.HttpResponseStatusCode, enp.FormatMsg(response.StatusCode))
	}
	if response.Body == nil {
		buf.WriteString("返回 response body nil\n")
		return nil, enp.Put(errorcode.HttpResponseBodyNil)
	}
	defer func(bc io.ReadCloser) {
		err = bc.Close()
		if err != nil {
			enp.Put(errorcode.ResponseBodyClose)
		}
	}(response.Body)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		buf.WriteString("返回 response body io ReadAll err \n")
		return nil, enp.Put(errorcode.IoReadAll, enp.AddError(err))
	}
	buf.WriteString(fmt.Sprintf("返回数据：%s\n", string(body)))
	buf.WriteString(fmt.Sprintf("执行时长: %d ms \n", end.Sub(start).Milliseconds()))
	return body, enp.Put(errorcode.Success)
}

// 系统接口测试参数封装
func AssembleUrl(address string) string {
	return fmt.Sprintf("%s%s", "http://web.trip-portal.qcs.dcsaas.cn", address)
}

func (param *RequestParam) AddHeader() {
	param.Header = map[string]interface{}{
		contentType: contentTypeForm,
	}
}

func (param *RequestParam) AddHeaderWithToken(token string) {
	param.Header = map[string]interface{}{
		contentType:      contentTypeForm,
		bo.AdminTokenKey: token,
	}
}

// 网约车连接参数封装
func AssembleTripUrl(address string) string {
	//return fmt.Sprintf("%s%s", "http://portal.qcs.dcsaas.cn", address) // prev
	//return fmt.Sprintf("%s%s", "http://portal.dev.dcsaas.cn", address) //dev
	return fmt.Sprintf("%s%s", config.Info().TripConfig.Host, address)
}

func (param *RequestParam) AddTripFormHeader() {
	param.Header = map[string]interface{}{
		contentType: contentTypeForm,
	}
}

func (param *RequestParam) AddTripJsonHeader() {
	param.Header = map[string]interface{}{
		contentType: contentTypeJSON,
	}
}
