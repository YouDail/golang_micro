package main

import (
	"flag"
	"fmt"
	"github.com/YouDail/golang_micro/hackathon-gateway/common"
	"github.com/YouDail/golang_micro/hackathon-gateway/handler"
	log "github.com/golang/glog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

//检查配置项
func VaildConf(s string) {

	if viper.GetString(s) == "" {
		log.Errorln("init", s, "  config is null")
		panic("启动失败！ 缺少配置参数:" + s)
	} else {
		log.Infoln("init 配置参数 ", s, "的值是: ", viper.GetString(s))
	}

}

func init() {

	log.Infoln("initial log config")
	_ = os.MkdirAll("logs", 0766)
	_ = flag.Set("alsologtostderr", "true")
	_ = flag.Set("stderrthreshold", "INFO")
	_ = flag.Set("log_dir", "logs")
	flag.Parse()

	defer log.Flush()

	curPath, err := os.Getwd()
	if err != nil {
		log.Errorln(err)
	}
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // or viper.SetConfigType("YAML")
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AddConfigPath(curPath)
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		log.Errorln(fmt.Errorf("Fatal error config file config.yaml : %s \n", err))
		panic("请检查配置文件config.yaml是否正确")
	}

	confs := []string{
		"registry.type",
		"registry.addr",
		"svc.Controller",
		"svc.Service",
	}

	for _, v := range confs {
		VaildConf(v)
	}

	//设置服务地址
	os.Setenv("MICRO_REGISTRY", viper.GetString("registry.type"))
	os.Setenv("MICRO_REGISTRY_ADDRESS", viper.GetString("registry.addr"))
	//设置注册服务超时时间
	os.Setenv("MICRO_CLIENT_RETRIES", "0")
	os.Setenv("MICRO_CLIENT_REQUEST_TIMEOUT", "10m")

	//设置服务监听地址
	f, err := os.OpenFile("service.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		log.Infoln("init  创建初始化配置文件service.json失败: ", err)
		panic(err)

	}

	len, err := f.WriteString(`{"listen":":"` + viper.GetString("httpPort") + `,"httpversion":1}`)
	if err != nil {
		log.Infoln("init  写入初始化配置文件service.json失败: ", err)
		panic(err)
	}
	f.Close()
	log.Infoln("init  写入初始化配置文件service.json成功，写入字节数:  ", len)

	os.Setenv("SERVICE_LISTEN", ":" + viper.GetString("httpPort"))

	log.Infoln("初始化redis连接")
	common.InitReCli()

	log.Infoln("RegisteMetrices to consul")
	common.RegisterMetrics()

}

func main() {
	app := iris.New()

	app.Get("/hackathon/maxClasses/{gradeId: int64 range(0,10)}", handler.HandleGradeId)

	log.Infoln("启用pprof")
	app.Use(pprof.New())

	log.Infoln("监听 tcp 0.0.0.0:" + viper.GetString("metrics.Port") + ", pprof地址：http://localhost:" + viper.GetString("metrics.Port") + "/debug/pprof")
	go func() {
		err := http.ListenAndServe(":"+viper.GetString("metrics.Port"), nil)
		if err != nil {
			log.Errorln("启动http server失败:", err)
		}
	}()

	//Prometheus client内置了golang metrics暴露的handler
	http.Handle("/metrics", promhttp.Handler())
	log.Infoln("Prometheus metrics 地址: http://localhost:" + viper.GetString("metrics.Port") + "/metrics")

	app.Run(iris.Addr(os.Getenv("SERVICE_LISTEN")))
}
