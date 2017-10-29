FROM alpine:3.6
MAINTAINER Utkarsh Mani Tripathi <utkarshmani1997@gmail.com>
ADD main /usr/local/bin
RUN apk add --no-cache \
     iproute2 \
     curl \
     net-tools \
     mii-tool \
     procps \
     libc6-compat

CMD main
EXPOSE 9500
