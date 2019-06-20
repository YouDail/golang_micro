package handler

import (
	"context"
	"errors"
	serviceProto "github.com/YouDail/golang_micro/hackathon-service/proto"
	log "github.com/golang/glog"
	"github.com/spf13/viper"
)

var hackathonSvc serviceProto.ClassInfoService

func HandleClassId(id int64) (string, error) {

	log.Infoln("create new hackathon client to the consul")

	hackathonSvc = GetClassClient(viper.GetString("registry.addr"), viper.GetString("svc.Service"))

	log.Infoln("call the hackathon")
	rsp, err := hackathonSvc.GetClassName(context.TODO(), &serviceProto.ClassIdRequest{
		ClassId: id,
	})

	if err != nil {
		log.Errorln(" rsp error: ", err)
		return "", errors.New(err.Error())
	}

	log.Infof("查询OK， 班级ID是%d, 班级名称是%s ", rsp.ClassId, rsp.ClassName)
	return rsp.ClassName, nil
}
