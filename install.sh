#!/bin/bash

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

go get github.com/golang/glog
go get github.com/golang/protobuf

go test -v github.com/qwb2333/Pierce/test
if [ $? -ne 0 ]; then
    echo 'unitests failed. check whether protobuf is already install and $GOPATH, $GOBIN, $GOROOT is correct.'
    exit
fi

go install main/inner.go
if [ $? -ne 0 ]; then
    exit
fi
go install main/outer.go

if [ ! -d "$GOBIN/../config" ]; then
  mkdir $GOBIN/../config
fi

cp $DIR/config/inner.properties $GOBIN/../config
cp $DIR/config/outer.properties $GOBIN/../config
cp $DIR/config/outer_dot.conf $GOBIN/../config

echo "install success."
echo "binary file are in $GOBIN"
echo "config file are in $GOBIN/../config"
echo -e "use command to run binary such as:\n"

echo 'cd $GOBIN'
echo -e "./inner ../config/inner.properties\n"

echo -e "or\n"

echo 'cd $GOBIN'
echo "./outer ../config/outer.properties ../config/outer_dot.conf"