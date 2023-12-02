#!/usr/bin/env bash

docker buildx build --platform linux/amd64,linux/arm64 -f sidecar.Dockerfile -t jitsucom/sidecar:latest --push .