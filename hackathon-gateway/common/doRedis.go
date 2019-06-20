package common

import (
	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"strconv"
)

var ReCli *redis.Client

func InitReCli() {
	glog.Infoln("实例化redis客户端")
	ReCli = initRedis()
}

//初始化redis客户端实例
func initRedis() *redis.Client {

	var client *redis.Client

	//获取库的编号
	var redisDB int
	redisDB, _ = strconv.Atoi(viper.GetString("redis.DB"))

	//默认库的编号为12
	if viper.GetString("redis.DB") == "" {
		client = redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis.Addr"),
			Password: "",
			DB:       12,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis.Addr"),
			Password: "",
			DB:       redisDB,
		})
	}

	pong, err := client.Ping().Result()
	if err != nil {
		//连接redis出错
		glog.Errorln("init connect to redis server error: ", err)
		panic(err)
	}
	glog.Infoln("init  result of ping redis server : ", pong)
	return client
}

//简单的K-V数据存储
type RedKV struct {
	RedKey string `json:"redKey"`
	RedVal string `json:"redVal"`
}

//set方法
func (s *RedKV) SetKV() (bool, error) {

	//默认 K-V 过期时间为60分钟
	err := ReCli.Set(s.RedKey, s.RedVal, 3600000000000).Err()
	if err != nil {
		return false, err
	}

	glog.Infoln("doRedis 成功插入键值对: ", s.RedKey, s.RedVal)

	return true, nil
}

//get方法
func (s *RedKV) GetKV() (bool, string) {

	val, err := ReCli.Get(s.RedKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, "key does not exist"
		}

		return false, err.Error()

	}
	glog.Infoln(" doRedis 查询value值是:", val)
	return true, val

}
