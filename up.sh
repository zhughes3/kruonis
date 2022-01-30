#!/bin/bash

docker-compose --env-file=cmd/timelinesv2/config.env up --build --detach
