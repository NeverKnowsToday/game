#!/bin/bash
#set -x

# 基础信息
CURRENT_DIR=$PWD
GATEWAY_NAME_PREFIX="game"
GATEWAY_IMAGETAG="latest"
GATEWAY_IMAGE_SERVER_NAME="game"
ALL_SERVICE_NAME="main"
DOCKERNETWORK="game"
SERVICETYPE="all"
DB_SERVICE_TYPE="container"
# server配置文件信息
SERVER_PATH="$CURRENT_DIR/server"
SERVER_SERVICE_NAME="game-server"
SERVER_YAML="server.yaml"

DB_YAML="mysql.yaml"
DB_PATH="$CURRENT_DIR/mysql"
DB_TEMPLATE_YAML="mysql.yaml"

function PrintHelp() {
  echo "Usage: "
  echo "  deploy.sh <mode>  [-s <service type>] "
  echo "    <mode> - one of 'up', 'down', 'restart',  'start', 'stop', 'pull' 'clearimages'"
  echo "      - 'up' - bring up the service(server, db) with docker-compose up"
  echo "      - 'down' - down the service and clear the database and configuration"
  echo "      - 'restart' - restart the service"
  echo "      - 'start'  - start the service"
  echo "      - 'stop'  - stop the service"
  echo "      - 'pull'  - pull service images"
  echo "      - 'clearimages'  - clear the service images"
  echo "      -s <service type> - the service type: server(server),db(mysql),all"
  echo "  deploy.sh -h (print this message)"
  echo
  echo "then bring up the network. e.g.:"
  echo "all defaults:"
  echo "	deploy.sh up    >>create game service and mysql "
  echo "	deploy.sh down  >>down fame service and clear database of mysql"
}

function ClearGATEWAYImages() {
  echo "**********clearing game images...**********"
  if [ "$SERVICETYPE" == "server" ]; then
    ClearServerImage
  elif [ "$SERVICETYPE" == "db" ]; then
    ClearDbImage
  elif [ "$SERVICETYPE" == "all" ]; then
    ClearServerImage
    ClearDbImage
  fi
}

function ClearServerImage() {
  echo "**********clearing SERVER images...**********"
  num=$(docker ps | grep "$GATEWAY_IMAGE_SERVER_NAME" | grep "$GATEWAY_IMAGETAG" | wc -l)
  if [ "$num" != "0" ]; then
    echo "please down the container of using $GATEWAY_IMAGE_SERVER_NAME:$GATEWAY_IMAGETAG"
  elif [ "$num" == "0" ]; then
    imageNum=$(docker images | grep "$GATEWAY_IMAGE_SERVER_NAME" | grep "$GATEWAY_IMAGETAG" | wc -l)
    if [ "$imageNum" != "0" ]; then
      docker rmi $GATEWAY_IMAGE_SERVER_NAME:$GATEWAY_IMAGETAG -f
    fi
  fi
}

function ClearDbImage() {
  echo "**********clearing db images...**********"
  num=$(docker ps | grep "mysql" | grep "5.7" | wc -l)
  if [ "$num" != "0" ]; then
    echo "please down the container of using mysql:5.7"
  elif [ "$num" == "0" ]; then
    imageNum=$(docker images | grep "mysql" | grep "5.7" | wc -l)
    if [ "$imageNum" != "0" ]; then
      docker rmi mysql:5.7 -f
    fi
  fi
}

function UpServer() {
  echo "***********up server*******************"
  cd $SERVER_PATH
  num=$(docker network ls | grep $DOCKERNETWORK | wc -l)
  if [ "$num" == "0" ]; then
    docker network create $DOCKERNETWORK
  fi
  docker-compose -f ./$SERVER_YAML up -d
  cd $CURRENT_DIR
}

function UpDB() {
  echo "***********up db*******************"
  if [ "${DB_SERVICE_TYPE}" == "container" ]; then
    cd $DB_PATH
    num=$(docker network ls | grep $DOCKERNETWORK | wc -l)
    if [ "$num" == "0" ]; then
      docker network create $DOCKERNETWORK
    fi
    docker-compose -f $DB_YAML up -d
    cd $CURRENT_DIR
  fi
}

function NetworkUp() {
  echo "***********network up*******************"

  if [ "${SERVICETYPE}" == "server" ]; then
    UpServer
  elif [ "${SERVICETYPE}" == "db" ]; then
    UpDB
  elif [ "${SERVICETYPE}" == "all" ]; then
    UpDB
    UpServer
  fi
}

function DownServer() {
  echo "***********down server*******************"
  cd $SERVER_PATH
  docker-compose -f $SERVER_YAML down -v
  cd $CURRENT_DIR
}

function DownDB() {
  echo "***********down db*******************"
  cd $DB_PATH
  docker-compose -f $DB_YAML down -v
  # clean all data
  rm -rf $DB_PATH/data
  rm -rf $SERVER_PATH/cert
  cd $CURRENT_DIR
}

function NetworkDown() {
  if [ "${SERVICETYPE}" == "server" ]; then
    DownServer
  elif [ "${SERVICETYPE}" == "db" ]; then
    DownDB
  elif [ "${SERVICETYPE}" == "all" ]; then
    echo "***********down all*******************"
    DownDB
    DownServer
    sudo -S rm -rf $CURRENT_DIR/mysql/data
    echo "log" > $CURRENT_DIR/server/fabric_gateway.log
  fi
  containerNum=$(docker ps -a | grep $DOCKERNETWORK | wc -l)
  networkNum=$(docker network ls | grep $DOCKERNETWORK | wc -l)
  if [ "$containerNum" == "0" -a "$networkNum" -ne "0" ]; then
    docker network rm $DOCKERNETWORK
  fi
}

function StopServer() {
  echo "***********stop server*******************"
  cd $SERVER_PATH
  if [ -f "./$SERVER_YAML" ]; then
    docker-compose -f ./$SERVER_YAML down
  fi
  cd $CURRENT_DIR
}

function StopDB() {
  echo "***********stop db*******************"
  cd $DB_PATH
  if [ -f "./$DB_YAML" ]; then
    docker-compose -f ./$DB_YAML down
  fi
  cd $CURRENT_DIR
}

function NetworkStop() {
  echo "***********network stop*******************"
  if [ "${SERVICETYPE}" == "server" ]; then
    StopServer
  elif [ "${SERVICETYPE}" == "db" ]; then
    StopDB
  elif [ "${SERVICETYPE}" == "all" ]; then
    StopDB
    StopServer
  fi
}

function StartServer() {
  echo "***********start server*******************"
  cd $SERVER_PATH
  if [ -f "./$SERVER_YAML" ]; then
    docker-compose -f ./$SERVER_YAML up -d
    docker logs -f "$SERVER_SERVICE_NAME"
  fi
  cd $CURRENT_DIR
}

function StartDB() {
  echo "***********start db*******************"
  cd $DB_PATH
  if [ -f "./$DB_YAML" ]; then
    docker-compose -f ./$DB_YAML up -d
  fi
  cd $CURRENT_DIR
}

function NetworkStart() {
  echo "***********network start*******************"
  if [ "${SERVICETYPE}" == "server" ]; then
    StartServer
  elif [ "${SERVICETYPE}" == "db" ]; then
    StartDB
  elif [ "${SERVICETYPE}" == "all" ]; then
    StartDB
    StartServer
  fi
}

function RestartServer() {
  echo "***********restart server*******************"
  cd $SERVER_PATH
  if [ -f "./$SERVER_YAML" ]; then
    docker-compose -f ./$SERVER_YAML down -v
    docker-compose -f ./$SERVER_YAML up -d
    docker logs -f "$SERVER_SERVICE_NAME"
  else
    echo "no configuration file"
    exit 1
  fi
}

function RestartDB() {
  echo "***********restart server*******************"
  cd $DB_PATH
  if [ -f "./$DB_YAML" ]; then
    docker-compose -f ./$DB_YAML down -v
    docker-compose -f ./$DB_YAML up -d
  else
    echo "no configuration file"
    exit 1
  fi

}

function NetworkRestart() {
  echo "***********network restart*******************"
  if [ "${SERVICETYPE}" == "server" ]; then
    RestartServer
  elif [ "${SERVICETYPE}" == "db" ]; then
    RestartDB
  elif [ "${SERVICETYPE}" == "all" ]; then
    RestartDB
    RestartServer
  fi
}

function PullServerImages() {
  echo "***********pull server image : $SERVICETYPE; tag: $GATEWAY_IMAGETAG*******************"
  docker pull $GATEWAY_IMAGE_SERVER_NAME:$GATEWAY_IMAGETAG
}

function PullDBImages() {
  echo "***********pull db image *******************"
  docker pull mysql:5.7
}

function PullImages() {
  if [ "${SERVICETYPE}" == "server" ]; then
    PullServerImages
  elif [ "${SERVICETYPE}" == "db" ]; then
    PullDBImages
  elif [ "${SERVICETYPE}" == "all" ]; then
    PullServerImages
    PullDBImages
  fi
}

MODE=$1
shift

while getopts "h?s:" opt; do
  case "$opt" in
  h | \?)
    PrintHelp
    exit 0
    ;;
  s) # -s #指定服务类型:SERVER, db, all
    SERVICETYPE=$OPTARG
    ;;
  esac
done

#Create the network using docker compose
if [ "${MODE}" == "up" ]; then
  NetworkUp
elif [ "${MODE}" == "down" ]; then ## Clear the network
  NetworkDown
elif [ "${MODE}" == "stop" ]; then ## stop the SERVER service
  NetworkStop
elif [ "${MODE}" == "start" ]; then ## start the SERVER service
  NetworkStart
elif [ "${MODE}" == "restart" ]; then ## restart the SERVER service
  NetworkRestart
elif [ "${MODE}" == "pull" ]; then ## pull SERVER and nginx images
  PullImages
elif [ "${MODE}" == "clearimages" ]; then ## clear the images' GATEWAY
  ClearGATEWAYImages
else
  PrintHelp
  exit 1
fi
