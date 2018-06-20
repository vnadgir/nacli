#!/bin/bash
docker rm -f $(docker ps -a -q)
docker run -d --name nats-0 -p 4222:4222 nats -p 4222
docker run -d --name nats-1 -p 14222:14222 nats -p 14222 --routes nats://localhost:4222 --cluster nats://localhost:4222
docker run -d --name nats-2 -p 24222:24222 nats -p 24222 --routes nats://localhost:4222 --routes nats://localhost:14222 --cluster nats://localhost:4222

docker run -d --link nats-0 --name stan nats-streaming -cid testing-stan -ns nats://nats-0:4222 -SDV

echo "Running nats servers...."
docker ps
