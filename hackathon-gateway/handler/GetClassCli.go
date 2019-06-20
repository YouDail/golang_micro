package handler

import (
	serviceProto "github.com/YouDail/golang_micro/hackathon-service/proto"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
)

func GetClassClient(registryName, serveName string) serviceProto.ClassInfoService {

	srv := grpc.NewService(
		micro.Name("Class.client"),
		micro.Registry(
			etcdv3.NewRegistry(func(options *registry.Options) {
				options.Addrs = []string{registryName}
			},
			),
		),
	)

	return serviceProto.NewClassInfoService(serveName, srv.Client())
}
