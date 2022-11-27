#! /bin/sh

CurrentPath=$PWD
ProjectPath=$(dirname $(dirname "$PWD"))
ServerPath=$(dirname $(dirname "$PWD"))/server
Tag=$1

cd $ProjectPath
tar zcvf server.tar.gz server && mv server.tar.gz $CurrentPath

cd $CurrentPath

docker build -t game:latest .
docker tag game:latest game:$Tag

rm -rf server.tar.gz

