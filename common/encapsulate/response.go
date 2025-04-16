package encapsulate

import (
	"encoding/json"
	"fmt"
	"github.com/project-template/errorcode"
	"runtime"
	"strings"
)

// Response 服务内使用，主要看一些错误发生的原因
type Response struct {
	// 这里是需要向外传递的
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []byte `json:"data,omitempty"`
	// 这里是内部记录错误使用的
	Err               string           `json:"err,omitempty"`
	In                string           `json:"in,omitempty"`
	Out               string           `json:"out,omitempty"`
	Message           string           `json:"inner_msg,omitempty"`
	MessageFormatArgs []any            `json:"message_format_args,omitempty"`
	ResponseStack     []*ResponseStack `json:"response_stack,omitempty"`
}

type ResponseStack struct {
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

// Reply 系统对外返回统一结构
type Reply struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (response *Response) Reply(vo interface{}) Reply {
	var reply = Reply{Code: response.Code, Msg: response.Msg}
	if strings.Contains(response.Msg, "#?#") {
		ms := strings.Split(response.Msg, "#?#")
		reply.Msg = ms[0]
	}
	if response.Data == nil || vo == nil {
		return reply
	}
	err := json.Unmarshal(response.Data, vo)
	if err != nil {
		fmt.Println("Reply Unmarshal err:", err)
		fmt.Println(string(response.Data))
		return Reply{Code: errorcode.JsonUnmarshal}
	}
	reply.Data = vo
	return reply
}

func (response *Response) withOption(opts ...Option) {
	for _, opt := range opts {
		opt.apply(response)
	}
}

type Option interface {
	apply(*Response)
}

type optionFunc func(*Response)

func (f optionFunc) apply(response *Response) {
	f(response)
}

func AddError(err error) Option {
	return optionFunc(func(response *Response) {
		response.Err = err.Error()
	})
}

func AddIn(args ...any) Option {
	return optionFunc(func(response *Response) {
		b, _ := json.Marshal(args)
		response.In = string(b)
	})
}

func FormatMsg(args ...any) Option {
	return optionFunc(func(response *Response) {
		response.MessageFormatArgs = append(response.MessageFormatArgs, args...)
	})
}

func AddOut(args ...any) Option {
	return optionFunc(func(response *Response) {
		b, _ := json.Marshal(args)
		response.Out = string(b)
	})
}

func AddData(args any) Option {
	return optionFunc(func(response *Response) {
		b, _ := json.Marshal(args)
		response.Data = b
	})
}

func Put(code int, opts ...Option) *Response {
	response := &Response{Code: code}
	response.withOption(opts...)
	if errorcode.ErrorCode == nil || len(errorcode.ErrorCode) == 0 {
		errorcode.AdminCode()
		errorcode.BaseCode()
		errorcode.MiddlewareCode()
		errorcode.RoleCode()
		errorcode.SdkCode()
		errorcode.LogCode()
	}
	response.Msg = errorcode.ErrorCode[code]
	response.Message = errorcode.ErrorCode[code]
	if response.MessageFormatArgs != nil && len(response.MessageFormatArgs) > 0 {
		response.Message = response.Msg
		if strings.Contains(response.Msg, "#?#") {
			ms := strings.Split(response.Msg, "#?#")
			response.Message = ms[1]
		}
		response.Message = fmt.Sprintf(response.Message, response.MessageFormatArgs...)
	}
	if code != errorcode.Success {
		stacks := make([]*ResponseStack, 0)
		for i := 1; i <= 5; i++ {
			_, file, line, _ := runtime.Caller(i)
			if len(file) == 0 || line == 0 {
				break
			}
			stacks = append(stacks, &ResponseStack{File: file, Line: line})
			fmt.Println(file, line)
		}
		response.ResponseStack = stacks
		b, _ := json.Marshal(response)
		fmt.Println(string(b))
	}
	return response
}
