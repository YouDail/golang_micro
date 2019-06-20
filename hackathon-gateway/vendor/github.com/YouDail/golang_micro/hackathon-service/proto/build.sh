rm -f class.micro.go  class.pb.go
protoc --proto_path=. --micro_out=. --go_out=.  class.proto
