#!/usr/bin/env bash
cd "$(dirname "$0")"

go run ./*.go localhost:9092 alpha 2 25