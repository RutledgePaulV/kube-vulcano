#!/usr/bin/env bash
command="docker run -v $PWD:/go/src/app -e PROJECT=kube-vulcano rutledgepaulv/godep-vendor-builder:master"
echo "Executing command to build go binary: ${command}"
${command}

command="docker build -t rutledgepaulv/kube-vulcano $PWD"
echo "Executing command to package binary in scratch image: ${command}"
${command}

command="docker push rutledgepaulv/kube-vulcano"
echo "Executing command to push new image to remote repository: ${command}"
${command}