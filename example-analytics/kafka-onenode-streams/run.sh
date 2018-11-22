#!/usr/bin/env bash
cd "$(dirname "$0")"

docker-compose up -d
printf "\n see http://localhost:8082 for a visualizer \n"