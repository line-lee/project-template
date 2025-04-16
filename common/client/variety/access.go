package variety

import (
	"encoding/json"
	"github.com/project-template/common/client/proto"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/errorcode"
)

const (
	// *****************************************************************
	// 特别注意：如果调整了顺序，需要同时更新服务的发现和注册
	// *****************************************************************
	adminModule = iota // admin相关路由
	roleModule         // 权限角色相关路由
)

func Do(moduleNum, apiNum int, data []byte) (*proto.Response, error) {
	var resp *enp.Response
	switch moduleNum {
	case adminModule:
		resp = adminFunctions[adminRout(apiNum)](data)
	case roleModule:
		resp = roleFunctions[roleRout(apiNum)](data)
	default:
		b, _ := json.Marshal(enp.Put(errorcode.UnknownModule, enp.AddIn(moduleNum, apiNum)))
		return &proto.Response{Data: b}, nil
	}
	b, err := json.Marshal(resp)
	if err != nil {
		b, _ = json.Marshal(enp.Put(errorcode.JsonMarshal, enp.AddIn(resp), enp.AddError(err)))
		return &proto.Response{Data: b}, nil
	}
	return &proto.Response{Data: b}, nil
}
