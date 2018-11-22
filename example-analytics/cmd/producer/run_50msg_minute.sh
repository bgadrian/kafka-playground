#!/usr/bin/env bash
cd "$(dirname "$0")"

#<broker> <topic> <partitionCount> <onlineUsersPeak> <eventsPerMinPerUser>
go run main.go localhost:9092 analytics 2 25 2