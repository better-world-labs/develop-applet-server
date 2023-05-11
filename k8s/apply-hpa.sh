#!/usr/bin/env sh

#usage: ./apply-hpa.sh $namespace $appName

cp hpa.yaml hpa.yaml.tmp

sed -i "" "s/__NS__/$1/" hpa.yaml.tmp
sed -i "" "s/__APP_NAME__/$2/" hpa.yaml.tmp

kubectl apply -f hpa.yaml.tmp