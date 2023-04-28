#!/usr/bin/env bash

err() { echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: [ERROR] $*" >&2; }

path=""
config=""
version=""
arguments=""

while [[ $# -gt 0 ]]; do
  case "$1" in
  -p | --path)
    path="$2"
    shift 2
    ;;
  -c | --config)
    config="$2"
    shift 2
    ;;
  -v | --version)
    version="$2"
    shift 2
    ;;
  *)
    arguments="$*"
    shift $# # past argument
    ;;
  esac
done


errors=""
if [ -z "${path}" ]; then errors="${errors}pass -p|--path; "; fi
if [ -z "${config}" ]; then errors="${errors}pass -c|--config; "; fi
if [ -z "${version}" ]; then errors="${errors}pass -v|--version; "; fi

arguments="${arguments} ${version}"

if [ -n "${errors}" ]; then err "${errors}" && exit 1; fi

images=$(docker images splitter | wc -l)
if [ "$images" -eq 1 ]; then
  make build-docker
fi

docker run --rm \
  -v ~/.splitter/config:/root/.splitter/config \
  -v ~/.ssh/id_rsa:/root/.ssh/id_rsa \
  -v "$path":"$path" \
  -v "$config":/tmp/config \
  splitter:latest \
  -c /tmp/config \
  "$arguments"
