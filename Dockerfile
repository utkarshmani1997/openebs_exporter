FROM alpine:3.6
MAINTAINER Utkarsh Mani Tripathi <utkarshmani1997@gmail.com>
ADD entrypoint.sh /usr/local/bin
ADD main /usr/local/bin
RUN chmod +x /usr/local/bin/entrypoint.sh
RUN apk add --no-cache libc6-compat  
ENTRYPOINT entrypoint.sh
EXPOSE 9500
