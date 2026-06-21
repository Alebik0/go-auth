#!/usr/bin/env bash
set -a

source ../.env

set +a

swag init

go run .
