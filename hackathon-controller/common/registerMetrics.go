package common

import (
	"encoding/json"
	log "github.com/golang/glog"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func RegisterMetrics() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Errorln(err)
		panic("GetIntranetIp err: " + err.Error())
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Infoln("ip:", ipnet.IP.String())
				err := RegisterToConsul(ipnet.IP.String())
				if err != nil {
					log.Errorln("注册metrics到consul失败! ")
				}
			}
		}
	}
}

func RegisterToConsul(ip string) error {
	log.Infoln("RegisterToConsul 接受参数: ", ip)
	var check CheckNode
	check.HTTP = "http://" + ip + ":" + viper.GetString("metrics.Port") + "/metrics"
	check.Interval = "5s"
	var metrics RegMetrics
	metrics.Name = viper.GetString("metrics.Name")
	metrics.Address = ip
	metrics.ID = viper.GetString("metrics.Name") + "-" + strings.Split(ip, ".")[3]
	metrics.MetricsPath = "/metrics"
	metrics.Port = viper.GetInt("metrics.Port")
	metrics.Tags = append(metrics.Tags, viper.GetString("metrics.Tag"), "go_process")
	metrics.Checks = append(metrics.Checks, check)
	metricsStr, err := json.Marshal(&metrics)
	if err != nil {
		log.Errorln("编码metrics失败: ", err)
		return err
	}

	//初始化http客户端
	client := &http.Client{}

	//将metrics注册到consul
	req, err := http.NewRequest("PUT", "http://"+viper.GetString("consulAddr")+"/v1/agent/service/register", strings.NewReader(string(metricsStr)))
	if err != nil {
		log.Errorln(" 将metrics注册到consul失败, http.NewRequest error:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorln("将metrics注册到consul失败, client.Do(req) error:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("将metrics注册到consul失败, ioutil.ReadAll(resp.Body) error:", err)

		return err
	}

	log.Infoln("成功将metrics注册到consul, resp.Body: ", string(body))
	return nil

}

type RegMetrics struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Port        int         `json:"port"`
	MetricsPath string      `json:"metrics_path"`
	Tags        []string    `json:"tags"`
	Checks      []CheckNode `json:"checks"`
}

type CheckNode struct {
	HTTP     string `json:"http"`
	Interval string `json:"interval"`
}
