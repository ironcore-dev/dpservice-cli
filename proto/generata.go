package dpdkproto

//go:generate  protoc --proto_path=. --go_out=../pkg/dpdkproto --go_opt=paths=source_relative --go-grpc_out=../pkg/dpdkproto --go-grpc_opt=paths=source_relative  dpdk.proto
