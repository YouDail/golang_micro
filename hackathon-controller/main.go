package main

import (
	"flag"
	"fmt"
	"github.com/YouDail/golang_micro/hackathon-controller/common"
	"github.com/YouDail/golang_micro/hackathon-controller/handler"
	proto "github.com/YouDail/golang_micro/hackathon-controller/proto"
	log "github.com/golang/glog"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
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
	viper.SetConfigName("key")  // name of config file (without extension)
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	viper.AddConfigPath(curPath)
	err = viper.MergeInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Errorln(fmt.Errorf("Fatal error config file key.yaml: %s \n", err))
		panic("请检查配置文件key.yaml是否正确")
	}

	confs := []string{
		"mysql.Addr",
		"mysql.User",
		"mysql.PasswdSecret",
		"mysql.DB",
		"registry.type",
		"registry.addr",
		"serviceName",
	}

	for _, v := range confs {
		VaildConf(v)
	}

	log.Infoln("RegisterMetrics to consul")
	common.RegisterMetrics()
}

func main() {

	log.Infoln("create n new service")

	os.Setenv("MICRO_REGISTRY", viper.GetString("registry.type"))
	os.Setenv("MICRO_REGISTRY_ADDRESS", viper.GetString("registry.addr"))
	os.Setenv("MICRO_SERVER_NAME", viper.GetString("serviceName"))

	//registre := etcdv3.NewRegistry()
	service := grpc.NewService(
		micro.Name(viper.GetString("serviceName")),
		micro.Registry(etcdv3.NewRegistry(func(options *registry.Options) {
			options.Addrs = []string{
				viper.GetString("registry.addr"),
			}
		})),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*60),
		micro.RegisterInterval(time.Second*60),
	)

	log.Infoln("init service")
	//service.Init()

	log.Infoln("register handler to etcdV3")
	err := proto.RegisterMaxClassesHandler(service.Server(), new(handler.GradeId))
	if err != nil {
		log.Errorln("register handler to etcdV3 error:", err)
		panic(err)
	}

	log.Infoln("初始化数据库引擎")
	common.InitDB()

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

	log.Infoln("run the server")
	err = service.Run()
	if err != nil {
		panic(err)
	}
}
