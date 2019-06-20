rm -f grade.micro.go	grade.pb.go 
protoc --proto_path=. --micro_out=. --go_out=.  grade.proto
