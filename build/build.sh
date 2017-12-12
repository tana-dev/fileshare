# resource file copy
cp -R ../src/resources ../bin/downloader/
cp ../src/config/user.json ../bin/downloader/config/

# compile
cd ../src/main/
GOOS=linux GOARCH=amd64 go build -o ../../bin/downloader/downloader main.go

