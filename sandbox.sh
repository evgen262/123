#!/usr/bin/env sh

APP_NAME=auth

eval $(minikube -p minikube docker-env)
docker build --no-cache --build-arg GOPROXY_LOGIN=$GOPROXY_LOGIN --build-arg GOPROXY_TOKEN=$GOPROXY_TOKEN --build-arg RELEASE=sandbox -t $APP_NAME:latest -f ./Dockerfile .
eval $(minikube docker-env -u)
helm upgrade --install -n sandbox $APP_NAME ./dev/helm -f ./dev/helm/helm-values-sandbox.yaml -f ./dev/helm/env-secret-sandbox.yaml
kubectl rollout restart deployment $APP_NAME-deployment -n sandbox
