package handler

import (
	"context"
	"github.com/YouDail/golang_micro/hackathon-service/common"
	proto "github.com/YouDail/golang_micro/hackathon-service/proto"
	log "github.com/golang/glog"
	"github.com/pkg/errors"
)

func (c *ClassId) GetClassName(ctx context.Context, req *proto.ClassIdRequest, rsp *proto.ClassNameResponse) error {

	log.Infoln("GetClassName client request: ", req.ClassId)

	//设置sql_mode
	setQry, err := common.DBengine.Query("set @@sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';")
	if err != nil {
		log.Errorln("GetClassName set @@sql_mode error: ", err)
	}
	log.Infoln("GetClassName set @@sql_mode  ok ", setQry)
	getSqlMode, err := common.DBengine.Query("SELECT @@SESSION.sql_mode as mode;")
	if err != nil {
		log.Errorln("GetClassName get @@SESSION.sql_mode error: ", err)
	}
	log.Infoln("GetClassName get @@SESSION.sql_mode ok, 结果是:", string(getSqlMode[0]["mode"]))

	//查询传参的班级名称
	classQry, err := common.DBengine.Query("SELECT name  from hackathon.hackathon_class WHERE id  = ?", req.ClassId)
	if err != nil {
		log.Errorln("GetClassName 查询班级name出错:", err)
		return err
	}
	log.Info("GetClassName 查询班级name结果是: ", classQry)

	if len(classQry) <= 0 {
		log.Errorln("GetClassName 查询班级name 为空，calss id不存在！:")
		return errors.New("calss id不存在！")
	}

	for _, v := range classQry {
		log.Infoln("GetClassName 查询的 班级名称： ", string(v["name"]))
		rsp.ClassId = req.ClassId
		rsp.ClassName = string(v["name"])
	}

	return nil
}
