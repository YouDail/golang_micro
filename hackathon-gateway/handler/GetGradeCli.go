package handler

import (
	controllerProto "github.com/YouDail/golang_micro/hackathon-controller/proto"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
)

func GetGradeClient(registryName, serveName string) controllerProto.MaxClassesService {

	srv := grpc.NewService(
		micro.Name("Grade.client"),
		micro.Registry(
			etcdv3.NewRegistry(func(options *registry.Options) {
				options.Addrs = []string{registryName}
			},
			)),
	)

	return controllerProto.NewMaxClassesService(serveName, srv.Client())
}
