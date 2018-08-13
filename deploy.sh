#!/usr/bin/env bash


export version=0.0.2
export containername=ethstat

echo "Test and Build Binary"
docker-compose -f docker-compose.yml run --rm unit
echo "Build image"
docker build -t $containername:$version .

echo "Pushing image to docker hub"
docker tag $containername:$version leondroid/$containername:$version
docker push leondroid/$containername:$version