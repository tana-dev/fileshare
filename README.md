# view file prepartion compile
$ cd ./src/main

$ go-bindata-assetfs ../resources/...

# compile cmd
$ cd ./src/main
## linux
$ GOOS=linux GOARCH=amd64 go build -o ../../bin/dl_tanaka-shu bindata_assetfs.go main.go

## windows
$ GOOS=windows GOARCH=amd64 go build -o ../../bin/dl_tanaka-shu.exe bindata_assetfs.go main.go

## mac
$ GOOS=darwin GOARCH=amd64 go build -o ../../bin/dl_tanaka-shu bindata_assetfs.go main.go

