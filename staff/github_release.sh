#!/bin/bash
# Example of use:
# TAG=v1.2 TOKEN=447879c5af5887eab22725605783e86d3304bc99 bash ./staff/github_release.sh
PROJECT="ezotrank/tslogs"
BIN_NAME="tslogs"
BIN_PATH="./tmp/$BIN_NAME"

if [ -z "$TAG" ]; then
  echo "TAG is empty please do 'export TAG=v0.1' for example"
  exit 1
fi

if [ -z "$TOKEN" ]; then
  echo "You forgot write yout TOKEN"
  exit 1
fi


gzip -9 -f $BIN_PATH || exit 1

response=`curl --data "{\\"tag_name\\": \\"$TAG\\",\\"target_commitish\\": \\"master\\",\\"name\\": \\"$TAG\\",\\"body\\": \\"Release of version $TAG\\", \\"draft\\": false,\\"prerelease\\": false}" \
  -H 'Accept-Encoding: gzip,deflate' --compressed "https://api.github.com/repos/$PROJECT/releases?access_token=$TOKEN" > response`

release_id=`cat response|head -n 10|grep '"id"'|head -n 1|awk '{print $2}'|sed -e 's/,//'`
rm response

if [ -z "$release_id" ]; then
  echo "something wrong"
  echo $response
  exit 1
fi

curl -X POST -H 'Content-Type: application/x-gzip' --data-binary @$BIN_PATH.gz "https://uploads.github.com/repos/$PROJECT/releases/$release_id/assets?name=$BIN_NAME.gz&access_token=$TOKEN"