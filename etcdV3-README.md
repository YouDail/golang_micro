
## etcdv3

#### pull镜像

```

docker login --username=mwcx998 registry.cn-hangzhou.aliyuncs.com
docker pull registry.cn-hangzhou.aliyuncs.com/ltd/etcd:v3.3.13
docker tag registry.cn-hangzhou.aliyuncs.com/ltd/etcd:v3.3.13 quay.io/coreos/etcd:v3.3.13 


```



#### 启动

```

docker run -d -p 22379:2379 -p 22380:2380 -d -v /data/etcdV3:/data quay.io/coreos/etcd:v3

```



#### 客户端调用API
curl -L http://localhost:2379/v3beta/kv/put -X POST -d '{"key": "Zm9v", "value": "YmFy"}'






### [Notes](https://github.com/etcd-io/etcd/blob/master/Documentation/dev-guide/api_grpc_gateway.md#notes)

gRPC gateway endpoint has changed since etcd v3.3:

```
* etcd v3.2 or before uses only [CLIENT-URL]/v3alpha/*.
* etcd v3.3 uses [CLIENT-URL]/v3beta/* while keeping [CLIENT-URL]/v3alpha/*.
* etcd v3.4 uses [CLIENT-URL]/v3/* while keeping [CLIENT-URL]/v3beta/*.
    * [CLIENT-URL]/v3alpha/* is deprecated.
* etcd v3.5 or later uses only [CLIENT-URL]/v3/*.
    * [CLIENT-URL]/v3beta/* is deprecated.
gRPC-gateway does not support authentication using TLS Common Name.
```

[Put and get keys](https://github.com/etcd-io/etcd/blob/master/Documentation/dev-guide/api_grpc_gateway.md#put-and-get-keys)

Use the /v3/kv/range and /v3/kv/put services to read and write keys:

```

curl -L http://localhost:2379/v3/kv/put \
  -X POST -d '{"key": "Zm9v", "value": "YmFy"}'
# {"header":{"cluster_id":"12585971608760269493","member_id":"13847567121247652255","revision":"2","raft_term":"3"}}

curl -L http://localhost:2379/v3/kv/range \
  -X POST -d '{"key": "Zm9v"}'
# {"header":{"cluster_id":"12585971608760269493","member_id":"13847567121247652255","revision":"2","raft_term":"3"},"kvs":[{"key":"Zm9v","create_revision":"2","mod_revision":"2","version":"1","value":"YmFy"}],"count":"1"}

# get all keys prefixed with "foo"
curl -L http://localhost:2379/v3/kv/range \
  -X POST -d '{"key": "Zm9v", "range_end": "Zm9w"}'
# {"header":{"cluster_id":"12585971608760269493","member_id":"13847567121247652255","revision":"2","raft_term":"3"},"kvs":[{"key":"Zm9v","create_revision":"2","mod_revision":"2","version":"1","value":"YmFy"}],"count":"1"}

[Watch keys](https://github.com/etcd-io/etcd/blob/master/Documentation/dev-guide/api_grpc_gateway.md#watch-keys)
Use the /v3/watch service to watch keys:
curl -N http://localhost:2379/v3/watch \
  -X POST -d '{"create_request": {"key":"Zm9v"} }' &
# {"result":{"header":{"cluster_id":"12585971608760269493","member_id":"13847567121247652255","revision":"1","raft_term":"2"},"created":true}}

curl -L http://localhost:2379/v3/kv/put \
  -X POST -d '{"key": "Zm9v", "value": "YmFy"}' >/dev/null 2>&1
# {"result":{"header":{"cluster_id":"12585971608760269493","member_id":"13847567121247652255","revision":"2","raft_term":"2"},"events":[{"kv":{"key":"Zm9v","create_revision":"2","mod_revision":"2","version":"1","value":"YmFy"}}]}}
```

[Transactions](https://github.com/etcd-io/etcd/blob/master/Documentation/dev-guide/api_grpc_gateway.md#transactions)


Issue a transaction with /v3/kv/txn:
```
# target CREATE
curl -L http://localhost:2379/v3/kv/txn \
  -X POST \
  -d '{"compare":[{"target":"CREATE","key":"Zm9v","createRevision":"2"}],"success":[{"requestPut":{"key":"Zm9v","value":"YmFy"}}]}'
# {"header":{"cluster_id":"12585971608760269493","member_id":"13847567121247652255","revision":"3","raft_term":"2"},"succeeded":true,"responses":[{"response_put":{"header":{"revision":"3"}}}]}
# target VERSION
curl -L http://localhost:2379/v3/kv/txn \
  -X POST \
  -d '{"compare":[{"version":"4","result":"EQUAL","target":"VERSION","key":"Zm9v"}],"success":[{"requestRange":{"key":"Zm9v"}}]}'
# {"header":{"cluster_id":"14841639068965178418","member_id":"10276657743932975437","revision":"6","raft_term":"3"},"succeeded":true,"responses":[{"response_range":{"header":{"revision":"6"},"kvs":[{"key":"Zm9v","create_revision":"2","mod_revision":"6","version":"4","value":"YmF6"}],"count":"1"}}]}
```

[Authentication](https://github.com/etcd-io/etcd/blob/master/Documentation/dev-guide/api_grpc_gateway.md#authentication)

Set up authentication with the /v3/auth service:

```
# create root user
curl -L http://localhost:2379/v3/auth/user/add \
  -X POST -d '{"name": "root", "password": "pass"}'
# {"header":{"cluster_id":"14841639068965178418","member_id":"10276657743932975437","revision":"1","raft_term":"2"}}

# create root role
curl -L http://localhost:2379/v3/auth/role/add \
  -X POST -d '{"name": "root"}'
# {"header":{"cluster_id":"14841639068965178418","member_id":"10276657743932975437","revision":"1","raft_term":"2"}}

# grant root role
curl -L http://localhost:2379/v3/auth/user/grant \
  -X POST -d '{"user": "root", "role": "root"}'
# {"header":{"cluster_id":"14841639068965178418","member_id":"10276657743932975437","revision":"1","raft_term":"2"}}

# enable auth
curl -L http://localhost:2379/v3/auth/enable -X POST -d '{}'
# {"header":{"cluster_id":"14841639068965178418","member_id":"10276657743932975437","revision":"1","raft_term":"2"}}
Authenticate with etcd for an authentication token using /v3/auth/authenticate:
# get the auth token for the root user
curl -L http://localhost:2379/v3/auth/authenticate \
  -X POST -d '{"name": "root", "password": "pass"}'
# {"header":{"cluster_id":"14841639068965178418","member_id":"10276657743932975437","revision":"1","raft_term":"2"},"token":"sssvIpwfnLAcWAQH.9"}
Set the Authorization header to the authentication token to fetch a key using authentication credentials:
curl -L http://localhost:2379/v3/kv/put \
  -H 'Authorization : sssvIpwfnLAcWAQH.9' \
  -X POST -d '{"key": "Zm9v", "value": "YmFy"}'
# {"header":{"cluster_id":"14841639068965178418","member_id":"10276657743932975437","revision":"2","raft_term":"2"}}

```

[etcd/Documentation/dev-guide/api_grpc_gateway.md](https://github.com/etcd-io/etcd/blob/master/Documentation/dev-guide/api_grpc_gateway.md)
