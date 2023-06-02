#!/usr/bin/env bash
docker build --progress plain --no-cache -f Dockerfile -t backend-test .
docker build --progress plain --no-cache -f Dockerfile-nginx -t tls-test .