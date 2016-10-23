from golang:alpine

run mkdir -p /go/src/app
workdir /go/src/app

copy ./cfg.go ./cfg/cfg.go
copy ./cfg_test.go ./cfg/cfg_test.go
run go test ./...
copy example/main.go ./main.go
run go-wrapper download
run go-wrapper install
cmd go-wrapper run