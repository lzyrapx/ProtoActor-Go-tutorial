# go install -v github.com/gogo/protobuf/protoc-gen-gogoslick
# https://github.com/AsynkronIT/protoactor-go/blob/dev/protobuf/protoc-gen-gograinv2/Makefile
# cp ~/go/bin/protoc-gen-gograinv2 $GOPATH/bin

protoc -I=. -I=$GOPATH/src --gogoslick_out=. protos.proto
protoc -I=. -I=$GOPATH/src --gograinv2_out=. protos.proto