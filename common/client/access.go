package client

import (
	"encoding/json"
	"github.com/project-template/common/client/proto"
	varietyapi "github.com/project-template/common/client/variety"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/errorcode"
)

func Do(service string, module, api int, data []byte) (*proto.Response, error) {
	if service == config.VarietyService {
		return varietyapi.Do(module, api, data)
	}
	b, _ := json.Marshal(enp.Put(errorcode.UnknownService, enp.AddIn(service)))
	return &proto.Response{Data: b}, nil
}
