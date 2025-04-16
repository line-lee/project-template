package proto

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/tools"
	"github.com/project-template/errorcode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"time"
)

var (
	varietyServiceConnection *grpc.ClientConn
	tripServiceConnection    *grpc.ClientConn
)

func GRPC(ctx context.Context, req *Request) *enp.Response {
	var conn *grpc.ClientConn
	var resp *enp.Response
	switch req.Service {
	case config.VarietyService:
		if varietyServiceConnection, resp = connection(varietyServiceConnection, config.Info().Core.Services[req.Service]); resp.Code != errorcode.Success {
			return resp
		}
		conn = varietyServiceConnection
	case config.TripService:
		if tripServiceConnection, resp = connection(tripServiceConnection, config.Info().Core.Services[req.Service]); resp.Code != errorcode.Success {
			return resp
		}
		conn = tripServiceConnection
	default:
		return enp.Put(errorcode.GrpcServiceUnknown, enp.AddIn(req.Service))
	}
	client := NewServiceClient(conn)
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	response, err := client.REQ(c, req)
	if err != nil {
		return enp.Put(errorcode.GrpcRequest, enp.AddIn(req), enp.AddError(err))
	}
	result := new(enp.Response)
	err = json.Unmarshal(response.Data, result)
	if err != nil {
		return enp.Put(errorcode.JsonUnmarshal, enp.AddIn(resp), enp.AddError(err))
	}
	return result
}

func connection(conn *grpc.ClientConn, register config.Register) (*grpc.ClientConn, *enp.Response) {
	var err error
	if conn == nil {
		url := fmt.Sprintf("dns:///%s:%d", register.Url, register.Grpc)
		fmt.Println("url:", url)
		option := grpc.WithTransportCredentials(insecure.NewCredentials())
		conn, err = grpc.NewClient(url, option)
		if err != nil {
			return nil, enp.Put(errorcode.GrpcNewClient, enp.AddIn(register), enp.AddError(err))
		}
	}
	return conn, enp.Put(errorcode.Success)
}

func GRPCStart(port int, server ServiceServer) {
	tools.SecureGo(func(args ...interface{}) {
		fmt.Println("===========================GRPCStart======================================")
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			fmt.Println("GRPCStart  net.Listen err -->", err)
			os.Exit(-1)
		}
		s := grpc.NewServer()
		RegisterServiceServer(s, server)
		fmt.Println("GRPCStart listening at", port)
		if err := s.Serve(lis); err != nil {
			fmt.Println("GRPCStart failed err : ", err)
		}
	})
}
