# set
OS=${1}

# make directory
mkdir -p ../bin/config/

# remove binary file
rm ../bin/macshare*

# resource file copy
cp -R ../src/resources ../bin/
cp ../src/config/user.json ../bin/config/

# compile
cd ../src/main/
case ${OS} in
    "linux" ) GOOS=linux GOARCH=amd64 go build -o ../../bin/macshare main.go
    ;;
    "windows" ) GOOS=windows GOARCH=amd64 go build -o ../../bin/macshare.exe main.go
    ;;
    "mac" ) GOOS=darwin GOARCH=amd64 go build -o ../../bin/macshare main.go
    ;;
    * ) echo "No compile"
    ;;
esac
