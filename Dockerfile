FROM gliderlabs/alpine:3.3

RUN apk add curl bash jq coreutils perl --no-cache

ADD assets/ /opt/resource
