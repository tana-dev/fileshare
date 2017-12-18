# resource file copy
cp -R ../src/resources ../bin/macshare/
cp ../src/config/user.json ../bin/macshare/config/

# compile
cd ../src/main/
GOOS=windows GOARCH=amd64 go build -o ../../bin/macshare/macshare.exe main.go
