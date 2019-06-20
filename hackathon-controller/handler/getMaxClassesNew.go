package handler

import (
	"context"
	"github.com/YouDail/golang_micro/hackathon-controller/common"
	proto "github.com/YouDail/golang_micro/hackathon-controller/proto"
	"github.com/go-xorm/xorm"
	log "github.com/golang/glog"
	"github.com/pkg/errors"
	"strconv"
)

//初始化DB引擎
var (
	DBengine *xorm.Engine
)

type GradeId struct{}

//根据年级编号查询其所有班级的男女生总人数并取出最大值
func (g *GradeId) GetMaxClassesNew(ctx context.Context, req *proto.GradeIdRequest, rsp *proto.GradeIdNewResponse) error {
	log.Infoln("GetMaxClass client request: ", req.GradeId)

	if req.GradeId == 0 {
		return errors.New("grade id不存在！")
	}

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

	//查询传参的年级名称
	gradeQry, err := common.DBengine.Query("SELECT name from hackathon.hackathon_grade WHERE id  = ?", req.GradeId)
	if err != nil {
		log.Errorf("GetMaxClass 查询年级编号为 %d 的 年级 name出错:", req.GradeId, err)
		return err
	}

	if len(gradeQry) <= 0 {
		log.Errorf("GetMaxClass 查询年级编号为 %d 的 年级 name 为空，grade id不存在！:", req.GradeId)
		return errors.New("grade id不存在！")
	}
	log.Infof("GetMaxClass 查询年级编号为 %d 的 年级 name 结果是: ", req.GradeId, gradeQry[0]["name"])

	//根据班级id查询班级id列表，然后查询班级的男女生总人数
	for _, v := range gradeQry {
		log.Infoln("GetMaxClass 循环 grade id值是: ", string(v["id"]))
		log.Infoln("GetMaxClass 循环 grade name值是 ", string(v["name"]))

		rsp.GradeId = req.GradeId
		rsp.GradeName = string(v["name"])

		//男生人数排序
		maleCountsQry, err := common.DBengine.Query("select * from (select classId,  sum(case when gender=1 then 1 else 0 end)  as male from  hackathon_student WHERE classId in ( SELECT id FROM hackathon_class WHERE gradeId = ? ) group by classId) e where male in (select male from (select male, count(distinct male) as count_male from (select * from (select classId,  sum(case when gender=1 then 1 else 0 end)  as male from  hackathon_student WHERE classId in ( SELECT id FROM hackathon_class WHERE gradeId = ? ) group by classId) a order by male desc) b group by classId order by male desc limit 1) c);", req.GradeId, req.GradeId)
		if err != nil {
			log.Error("GetMaxClass 统计各个班级的男生人数出错:", err)
			return err
		}

		for _, x := range maleCountsQry {
			log.Infof("%d 班级男生人数的班级，classId 是 %s， 其男生人数是 %s", req.GradeId, x["classId"], x["male"])

			maleCount, err := strconv.ParseInt(string(x["male"]), 10, 64)
			if err != nil {
				println("字符串转int64失败", err)
				return err
			}

			classId, err := strconv.ParseInt(string(x["classId"]), 10, 64)
			if err != nil {
				println("字符串转int64失败", err)
				return err
			}

			log.Infoln("组装返回数据rsp.Male ")
			var per proto.MaxNode
			per.Count = maleCount
			per.ClassId = append(per.ClassId, classId)
			rsp.Male = &per

		}

		//女生人数排序
		femaleCountsQry, err := common.DBengine.Query("select * from (select classId,  sum(case when gender=2 then 1 else 0 end)  as female from  hackathon_student WHERE classId in ( SELECT id FROM hackathon_class WHERE gradeId = ? ) group by classId) e where female in (select female from (select female, count(distinct female) as count_female from (select * from (select classId,  sum(case when gender=2 then 1 else 0 end)  as female from  hackathon_student WHERE classId in ( SELECT id FROM hackathon_class WHERE gradeId = ? ) group by classId) a order by female desc) b group by classId order by female desc limit 1) c);", req.GradeId, req.GradeId)
		if err != nil {
			log.Error("GetMaxClass 统计各个班级的女生人数出错:", err)
			return err
		}

		for _, x := range femaleCountsQry {
			log.Infof("%d 班级女生人数的班级，classId 是 %s， 其女生人数是 %s", req.GradeId, x["classId"], x["female"])

			femaleCount, err := strconv.ParseInt(string(x["female"]), 10, 64)
			if err != nil {
				println("字符串转int64失败", err)
				return err
			}

			classId, err := strconv.ParseInt(string(x["classId"]), 10, 64)
			if err != nil {
				println("字符串转int64失败", err)
				return err
			}

			log.Infoln("组装返回数据rsp.Female ")
			var per proto.MaxNode
			per.Count = femaleCount
			per.ClassId = append(per.ClassId, classId)
			rsp.Female = &per

		}

	}

	return nil
}

/*


//统计编号为2的年级的所有班级各自的男生人数，去重后取出男生人数总计的最大值；
select male from (select male, count(distinct male) as count_male from (select * from (select classId,  sum(case when gender=1 then 1 else 0 end)  as male from  hackathon_student WHERE classId in ( SELECT id FROM hackathon_class WHERE gradeId = 2 ) group by classId) a order by male desc) b group by classId order by male desc limit 1) b;

//根据上一步去重的结果，找出male总数符合去重结果的classId
select * from (select classId,  sum(case when gender=1 then 1 else 0 end)  as male from  hackathon_student WHERE classId in ( SELECT id FROM hackathon_class WHERE gradeId = 2 ) group by classId) e where male in (select male from (select male, count(distinct male) as count_male from (select * from (select classId,  sum(case when gender=1 then 1 else 0 end)  as male from  hackathon_student WHERE classId in ( SELECT id FROM hackathon_class WHERE gradeId = 2 ) group by classId) a order by male desc) b group by classId order by male desc limit 1) c);

*/
