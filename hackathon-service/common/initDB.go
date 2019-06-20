package common

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	log "github.com/golang/glog"
	"github.com/spf13/viper"
	"time"
)

//初始化DB引擎
var (
	DBengine *xorm.Engine
)

func GetDB(conStr string) *xorm.Engine {
	log.Infoln("init  mysql连接信息：", conStr)

	engine, err := xorm.NewEngine("mysql", conStr) //创建mysql连接
	if err != nil {
		log.Errorln("init  connect mysql db error: ", err)
	}

	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai") //设置时区
	engine.ShowSQL(true)                                      //会在控制台打印SQL语句
	//engine.Logger().SetLevel(core.LOG_DEBUG)                  //日志级别是debug
	engine.SetMapper(core.GonicMapper{}) //设置名称映射规则是结构体名称和对应的表名称以及结构体field名称与对应的表字段名称相同的命名；

	tabMet, _ := engine.DBMetas()
	log.Infoln("init  获取表结构信息", tabMet)

	//设置最大连接数
	setMaxCon, err := engine.Query("SET GLOBAL max_connections = 3000;")
	if err != nil {
		log.Errorln("initDB set GLOBAL max_connections error: ", err)
	}
	log.Infoln("GetClassName set GLOBAL max_connections  ok ", setMaxCon)
	getMaxCon, err := engine.Query("SELECT  @@GLOBAL.max_connections  as conn;")
	if err != nil {
		log.Errorln("GetClassName get @@GLOBAL.max_connections error: ", err)
	}
	log.Infoln("GetClassName SELECT  @@GLOBAL.max_connections; ok, 结果是:", string(getMaxCon[0]["conn"]))

	return engine
}

//初始化数据库连接
func InitDB() {

	var conStr string

	status, decryptData := Decrypt(viper.GetString("mysql.PasswdSecret"))
	if status {
		log.Infoln("解密mysql.PasswdSecret结果是:", decryptData)

	} else {
		panic("无法解密mysql.PasswdSecret " + decryptData)
	}

	conStr = viper.GetString("mysql.User") + ":" + decryptData + "@tcp(" + viper.GetString("mysql.Addr") + ")/" + viper.GetString("mysql.DB") + "?charset=utf8mb4"

	//conStr = "user:passwd@tcp(10.52.26.3:3306)/hackathon?charset=utf8mb4"

	DBengine = GetDB(conStr)

}
