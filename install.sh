#!/usr/bin/env bash

sudo apt-get install librdkafka-dev=0.11.3-1build1
docker-compose -f ./cluster-onenode/docker-compose.yml pull

echo "Install sudo apt-get install librdkafka-dev=0.11.6~1confluent5.0.1-1 from the Confluence repo https://docs.confluent.io/current/installation/installing_cp/deb-ubuntu.html#get-the-software"

