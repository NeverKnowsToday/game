#!/bin/bash

CURRENT_DIR=$PWD
BUILD_BIN=$PWD/build_bin
BUILD_IMAGE=$PWD/build_image

Tag=$1

function print_help() {
  echo "Usage: "
  echo "  build.sh <mode> "
  echo "    <mode> - one of 'bin', 'image'"
  echo "      - 'bin' - build the server bin in docker"
  echo "      - 'image' - down the service and clear the database and configuration"
  echo "    -t <image tag> - the fabric_gateway version(eg: v1.0 ...)"
}

function build_image() {
  if [ $Tag == "" ]; then
    Tag=v1.0
  fi
  cd $BUILD_IMAGE
  ./build.sh $Tag
  cd $CURRENT_DIR
}

function build_bin() {
  cd $BUILD_BIN
  ./build.sh
  cd $CURRENT_DIR
}

MODE=$1
shift

while getopts "h?:t:" opt; do
  case "$opt" in
  h | \?)
    PrintHelp
    exit 0
    ;;
  t) # -t 指定编译的镜像的tag
    Tag=$OPTARG
    ;;
  esac
done

if [ "${MODE}" == "image" ]; then
  build_image
elif [ "${MODE}" == "bin" ]; then
  build_bin
else
  print_help
  exit 1
fi
