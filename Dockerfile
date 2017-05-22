FROM gliderlabs/alpine:3.3

RUN apk add curl bash jq coreutils --no-cache

ADD assets/ /opt/resource
