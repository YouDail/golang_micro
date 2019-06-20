## hackathon 

A small go micro project



#### run etcd3

```
mkdir /data/etcdV3
chmod a+rw -R /data/etcdV3

docker run -d --restart=always  -p 22379:2379 -p 22380:2380 -d -v /data/etcdV3:/data quay.io/coreos/etcd:v3

```


#### run redis

```

docker run --name myredis --network=host -d redis

```

#### run mysql

```

mkdir /usr/mysqldata
chmod a+w -R /usr/mysqldata
docker run --name mymysql --network=host -v /usr/mysqldata:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7

导入数据
mysql -uroot -p123456 -h 127.0.0.1 -p 3306
source hackathon_class.sql;
source hackathon_grade.sql;
source hackathon_student.sql;


```

#### run controller & service

```

cd hackathon-controller
make docker
docker run -d --restart=always  --network=host hackathon-controller


cd hackathon-service
make docker
docker run -d --restart=always  --network=host hackathon-service

```


#### run gateway


```


cd hackathon-gateway
make docker
docker run -d --restart=always -m 4096M --memory-swap=4096M   --network=host hackathon-gateway


```


#### try it 


```

GET http://127.0.0.1:10010/hackathon/maxClasses/2

```


#### benchmark

```

cd hackathon-gateway/test
go test -v -bench=".*"

```

#### 使用vegeta进行压测
```
wget https://github.com/tsenart/vegeta/releases/download/cli%2Fv12.5.1/vegeta-12.5.1-linux-amd64.tar.gz
tar zxf vegeta-12.5.1-linux-amd64.tar.gz
mv vegeta /usr/local/bin/

先来个1600个
for m in {1..100};do for i in {0..15};do echo -e "{\"method\":\"GET\",\"url\":\"http://10.52.26.3:10010/hackathon/maxClasses/$i\"}";done;done |  vegeta  -profile heap attack -lazy -timeout 5s -workers 100 -format=json | tee results.bin | vegeta report  -type="hist[0,10ms,20ms,30ms]"
再来个1.6万个
for m in {1..1000};do for i in {0..15};do echo -e "{\"method\":\"GET\",\"url\":\"http://10.52.26.3:10010/hackathon/maxClasses/$i\"}";done;done |  vegeta  -profile heap attack -lazy -timeout 5s -workers 100 -format=json | tee results.bin | vegeta report  -type="hist[0,5ms,10ms,15ms,20ms,30ms]"
再压个16万个
for m in {1..1000};do for i in {0..15};do echo -e "{\"method\":\"GET\",\"url\":\"http://10.52.26.3:10010/hackathon/maxClasses/$i\"}";done;done |  vegeta  -profile heap attack -lazy -timeout 2s -workers 1000 -format=json | tee results.bin | vegeta report  -type=json
```

#### Go tool pprof & dstat

```

centos7安装graphviz包
yum install http://rpmfind.net/linux/centos/7.6.1810/os/x86_64/Packages/graphviz-2.30.1-21.el7.x86_64.rpm

查看gateway的runtime CPU和内存性能
while true;do go tool pprof --text  http://10.52.26.3:10909/debug/pprof/profile;go tool pprof --text  http://10.52.26.3:10909/debug/pprof/heap;sleep 5;done


查看系统进程级别的内存、cpu、io情况
yum install dstat -y
dstat --top-cpu --top-io --top-mem


```


#### Prometheus  metrics

```

curl http://127.0.0.1:10909/metrics

```

#### run consul：
```
docker run -d -p 8500:8500 consul
```


#### prometheus配置文件

```
# cat prometheus.yml
global:
  scrape_interval: 5s
  scrape_timeout: 5s
  evaluation_interval: 15s
scrape_configs:
  - job_name: hackathon
    relabel_configs:
    - source_labels:  ["__meta_consul_service"]
      regex: "(.*)"
      replacement: $1
      action: replace
      target_label: ""
    - source_labels:  ["__meta_consul_tags"]
      regex: "(.*),(.*),(.*)"
      replacement: $2
      action: replace
      target_label: "hackathon"
    metrics_path: /metrics
    scheme: http
    consul_sd_configs:
      - server: 10.52.26.3:8500
        scheme: http
```

#### run prometheus & granafa
```

docker run -d -v /usr/local/prometheus_my/prometheus.yml:/etc/prometheus/prometheus.yml -p 19090:9090   prom/prometheus

docker run -d -v /usr/local/grafana_plugins/:/var/lib/grafana/plugins -p 3000:3000 grafana/grafana


```

dashboards推荐[240](https://grafana.com/dashboards/240)


#### something perfect


Do you seen some sql-script in my project ?



```

//从mysql5.7以后，默认的sql_mode不允许对别名字段进行group by，这里需要去掉ONLY_FULL_GROUP_BY的mode设置
//下面这行设置仅对新款库生效

set @@sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
//下面这行设置对现有库生效
set sql_mode ='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';

//统计年级id是2的所有班级的各自的男女生总人数
select classId,  sum(case when gender=1 then 1 else 0 end)  as maleCount, sum(case when gender=2 then 1 else 0 end) as female from  hackathon_student WHERE classId = 2;

//统计年级id是2的所有班级的男生人数最多的班级
select * from (select classId,  sum(case when gender=1 then 1 else 0 end)  as male from  hackathon_student WHERE classId in ( SELECT id FROM hackathon_class WHERE gradeId = 2 ) group by classId) a order by male desc limit 1;

//统计年级id是2的所有班级的女生人数最多的班级
select * from (select classId,  sum(case when gender=2 then 1 else 0 end) as female from  hackathon_student WHERE classId in ( SELECT id FROM hackathon_class WHERE gradeId = 2 ) group by classId) a order by female desc limit 1;



```


#### dev stack

[go-micro](https://github.com/micro/go-micro)

[protobuf](https://github.com/golang/protobuf)

[iris](https://github.com/kataras/iris)

[xorm](https://github.com/go-xorm/xorm)

[crypto/aes](https://golang.google.cn/pkg/crypto/aes/)

[etcd3](https://github.com/micro/go-plugins/tree/master/registry/etcdv3)

[consul](https://www.consul.io/docs/agent/services.html)

[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)

[go-redis/redis](https://github.com/go-redis/redis)

[viper](https://github.com/spf13/viper)

[vegeta](https://github.com/tsenart/vegeta)

[golang-tool-pprof](https://github.com/iotd/jackliu-go-programming-note/blob/master/Golang-tool-pprof.md)

[prometheus/client_golang](https://github.com/prometheus/client_golang/tree/master/prometheus/promhttp)

[grafana](https://grafana.com/dashboards/240)


[iris/jwt](https://studyiris.com/example/exper/jwt.html)  待完善




Have fun !
