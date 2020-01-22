#!/usr/bin/bash

APP=KubeTestPod

rm -f ${APP} 
CGO_ENABLED=0 go build -o ${APP}
