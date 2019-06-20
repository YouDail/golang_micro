package handler

import (
	"context"
	"github.com/YouDail/golang_micro/hackathon-service/common"
	proto "github.com/YouDail/golang_micro/hackathon-service/proto"
	log "github.com/golang/glog"
	"github.com/pkg/errors"
	"strconv"
)

type ClassId struct{}

func (c *ClassId) GetClassInfoNew(ctx context.Context, req *proto.ClassIdRequest, rsp *proto.ClassIdNewResponse) error {
	log.Infoln("GetMaxClass client request: ", req.ClassId)

	//设置sql_mode
	setQry, err := common.DBengine.Query("set @@sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';")
	if err != nil {
		log.Errorln("GetMaxClass set @@sql_mode error: ", err)
	}
	log.Infoln("GetMaxClass set @@sql_mode  ok ", setQry)
	getSqlMode, err := common.DBengine.Query("SELECT @@SESSION.sql_mode as mode;")
	if err != nil {
		log.Errorln("GetMaxClass get @@SESSION.sql_mode error: ", err)
	}
	log.Infoln("GetMaxClass get @@SESSION.sql_mode ok, 结果是:", string(getSqlMode[0]["mode"]))

	//查询传参的班级名称
	classQry, err := common.DBengine.Query("SELECT name, gradeId from hackathon.hackathon_class WHERE id  = ?", req.ClassId)
	if err != nil {
		log.Errorln("GetClass 查询班级name出错:", err)
		return err
	}
	log.Info("GetClass 查询班级name结果是: ", classQry)

	if len(classQry) <= 0 {
		log.Errorln("GetClass 查询班级name 为空，calss id不存在！:")
		return errors.New("calss id不存在！")
	}

	for _, v := range classQry {
		log.Infoln("GetClass 查询的 班级名称： ", string(v["name"]))
		rsp.ClassId = req.ClassId
		rsp.ClassName = string(v["name"])
		log.Infoln("GetClass 查询的班级所属的年级id： ", string(v["gradeId"]))

		log.Infoln("GetClass 查询班级的男女生人数")
		studentsCountsQry, err := common.DBengine.Query("select classId,  sum(case when gender=1 then 1 else 0 end)  as maleCount, sum(case when gender=2 then 1 else 0 end) as femaleCount from  hackathon_student WHERE classId = ?;", req.ClassId)
		if err != nil {
			log.Errorln("GetClass 查询班级男女生人数情况出错:", err)
			return err
		}

		for _, x := range studentsCountsQry {
			maleCount, _ := strconv.ParseInt(string(x["maleCount"]), 10, 64)
			log.Infoln("GetClass 班级男生人数是:", maleCount)
			var person proto.Person
			person.MaleCount = maleCount

			femaleCount, _ := strconv.ParseInt(string(x["femaleCount"]), 10, 64)
			log.Infoln("GetClass 班级女生人数是:", femaleCount)
			person.FemaleCount = femaleCount

			rsp.Counts = append(rsp.Counts, &person)
		}

	}

	return nil

}
