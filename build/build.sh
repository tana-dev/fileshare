# set
OS=${1}

# resource file copy
cp -R ../src/resources ../bin/macshare/
cp ../src/config/user.json ../bin/macshare/config/

# remove binary file
rm ../bin/macshare/macshare*

# compile
cd ../src/main/
case ${OS} in
    "linux" ) GOOS=linux GOARCH=amd64 go build -o ../../bin/macshare/macshare main.go
    ;;
    "windows" ) GOOS=windows GOARCH=amd64 go build -o ../../bin/macshare/macshare.exe main.go
    ;;
    "mac" ) GOOS=darwin GOARCH=amd64 go build -o ../../bin/macshare/macshare main.go
    ;;
    * ) echo "No compile"
    ;;
esac
