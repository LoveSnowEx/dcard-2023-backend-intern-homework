#!/bin/bash

set -e

docker build --platform linux/amd64 -t registry.fly.io/page-list-service:latest .
flyctl auth docker
docker push registry.fly.io/page-list-service:latest
flyctl deploy --app page-list-service --image registry.fly.io/page-list-service:latest --remote-only
