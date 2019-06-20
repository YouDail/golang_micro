package handler

import (
	"context"
	"encoding/json"
	controllerProto "github.com/YouDail/golang_micro/hackathon-controller/proto"
	"github.com/YouDail/golang_micro/hackathon-gateway/common"
	log "github.com/golang/glog"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"strconv"
)

var hackathonGrd controllerProto.MaxClassesService

func HandleGradeId(ctx iris.Context) {
	log.Infoln("HandleGradeId 客户端地址: ", ctx.RemoteAddr())
	log.Infoln("HandleGradeId 客户端header: ", ctx.GetHeader("orgin"))

	gradeId, err := ctx.Params().GetInt64("gradeId")
	if err != nil {
		log.Errorln("HandleGradeId get gradeId Params error:", err)
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"code":    400,
			"message": err.Error(),
		})

		return
	}
	log.Infoln("客户端传参GradeId: ", gradeId)

	if gradeId == 0 || gradeId == 3 || gradeId > 7 {

		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString(`{"projectName":"hackathon","teamName":"海豹队","teamMembers":["孙士才","张肖辉","华玉磊","李腾达"],"finishTime":"2019-06-12 11:50:00"}`)

		return
	}

	log.Infoln("准备从redis取出数据")
	var kv common.RedKV
	kv.RedKey = "/hackathon/maxClasses/" + strconv.FormatInt(gradeId, 10)
	log.Infoln("预备 key: ", kv.RedKey)
	stauts, resDat := kv.GetKV()
	if stauts {
		log.Infoln("成功从redis里取出数据，下面开始校验数据")

		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString(string(resDat))
		return

	} else {
		log.Errorln("从redis取数据失败! ", resDat)
	}

	log.Infoln("create new hackathon client to the consul")
	hackathonGrd = GetGradeClient(viper.GetString("registry.addr"), viper.GetString("svc.Controller"))
	if hackathonGrd == nil {
		log.Errorln("连接controller rpc服务失败!")
	}

	log.Infoln("call the hackathon")
	rsp, err := hackathonGrd.GetMaxClassesNew(context.TODO(), &controllerProto.GradeIdRequest{
		GradeId: gradeId,
	})

	if err != nil {
		log.Errorln(" rsp error: ", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": common.GetErr(err.Error()),
		})

		return
	}

	log.Infoln("rsp: ", rsp)
	log.Infof("查询OK， 年级ID是%d, 年级名称是%s, 男生人数最多是 %d, 男生人数最多的班级id是%d, 女生人数最多是 %d, 女生人数最多的班级id是%d ", rsp.GradeId, rsp.GradeName, rsp.Male.Count, rsp.Male.ClassId, rsp.Female.Count, rsp.Female.ClassId)

	var (
		myRes   Resp
		myclass ClassNode
	)

	//循环男生人数最多的班级
	for _, x := range rsp.Male.ClassId {
		className, err := HandleClassId(x)
		if err != nil {
			log.Error("查询班级名称失败！", err)
			return
		}

		myclass.ClassName = className
		myclass.ClassId = x

		myRes.MaleClasses = append(myRes.MaleClasses, myclass)
	}

	//循环女生人数最多的班级
	for _, x := range rsp.Female.ClassId {
		className, err := HandleClassId(x)
		if err != nil {
			log.Error("查询班级名称失败！", err)
			return
		}

		myclass.ClassName = className
		myclass.ClassId = x

		myRes.FemaleClasses = append(myRes.FemaleClasses, myclass)
	}

	myRes.GradeName = rsp.GradeName
	myRes.Male = rsp.Male.Count
	myRes.Female = rsp.Female.Count

	aabbc, _ := json.Marshal(&myRes)
	log.Infoln("返回数据 STRING:", string(aabbc))

	log.Infoln("准备将数据存入redis")
	kv.RedVal = string(aabbc)
	status, err := kv.SetKV()
	if status {
		log.Infoln("成功将数据存入redis")
	} else {
		log.Errorln("将数据存入redis失败: ", err)
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.WriteString(string(aabbc))
	return
}
