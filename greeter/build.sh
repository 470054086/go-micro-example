#!/usr/bin/env bash



###创建micro-service/api
cd /home/www/api/
docker build -t micro-service/api -f ./Dockerfile_api .
###创建micro-service/api-service
docker build -t micro-service/api-service  .

###创建micro-service/order-service
cd /home/www/srvtwo

docker build -t micro-service/order-service  .

###创建micro-service/say-service

cd /home/www/srv

docker build -t micro-service/say-service  .
###运行docker-compose
cd /home/www

docker-compose up -d